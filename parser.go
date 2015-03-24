package ghostdoc

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/npolar/ciface"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	textFileRegex = `(?i)^.+\.txt$`
	jsonFileRegex = `(?i)^.+\.json|geojson|topojson$`
	csvFileRegex  = `(?i)^.+\.csv|tsv|txt`
	newlineRegex  = `\n|\r\n|\n\r$`
)

type Parser struct {
	Cli         *cli.Context
	DataChannel chan interface{}
	Input       []string
}

// NewParser initializes the Parser struct and passes a pointer to the caller
func NewParser(c *cli.Context, dc chan interface{}) *Parser {
	return &Parser{
		Cli:         c,
		DataChannel: dc,
		Input:       c.Args(), // Load the commandline arguments
	}
}

// Parse the provided input
func (p *Parser) Parse() {
	p.checkInput()
}

// checkInput performs basic argument checks and throws an error or triggers an input handler
func (p *Parser) checkInput() {
	if len(p.Input) > 0 && len(p.Input[0]) > 0 { // Check for at least one argument
		p.handleInput()
	} else {
		name := p.Cli.App.Name
		log.Fatalln(name, "called without an argument. See", name, "-h for usage information.")
	}
}

// @TODO go routine here?
// handleInput triggers file or directory parsing depending on the input type
func (p *Parser) handleInput() {
	for _, input := range p.Input {
		state, err := os.Stat(input)
		p.handleFatal(err)

		if state.IsDir() {
			p.globDir(input)
		} else {
			p.parseFile(input)
		}
	}
}

// globDir reads all compatible files in the directory and calls parseFile on each of them
func (p *Parser) globDir(dir string) {
	path, err := filepath.Abs(dir)       // Resolve input to an absolute path
	entries, err := ioutil.ReadDir(path) // Read all files in the directory
	p.handleFatal(err)                   // Check for errors

	for _, entry := range entries {
		if p.Cli.Bool("recursive") && entry.IsDir() {
			p.globDir(entry.Name())
		} else {
			if fname := entry.Name(); p.parsable(fname) {
				p.parseFile(fname)
			}
		}
	}
}

// parseFile parses the parser input (p.Input)
func (p *Parser) parseFile(fname string) {
	raw, err := ioutil.ReadFile(fname)

	switch {
	case p.isText(fname):
		p.parseText(raw, fname)
	case p.isJson(fname):
		err = p.parseJson(raw, fname)
	case p.isCsv(fname):
		err = p.parseCsv(raw, fname)
	default:
		log.Println("Unsupported file type")
	}

	if err != nil {
		log.Println(err)
	}
}

// parseText reads a text file removes any new lines and extra white spaces
// and turns it into json using the payload key. the resulting interface is
// pushed onto the DataChannel for further processing
func (p *Parser) parseText(data []byte, fname string) error {
	var textMap = make(map[string]interface{})

	text := p.replaceNewLines(data, " ")
	textMap[p.Cli.String("payload-key")] = strings.TrimSpace(string(text))

	err := p.parseFileName(fname, &textMap)

	p.DataChannel <- textMap

	return err
}

// parseJson reads json data into an interface and pushes the result onto
// the data channel for further processing
func (p *Parser) parseJson(data []byte, fname string) error {
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)

	if err == nil {
		err = p.parseFileName(fname, &jsonData)

		p.DataChannel <- jsonData
	}

	return err
}

// parseCsv reads  csv data into a []interface{} using github.com/npolar/ciface.
// Each interface{} is then pushed onto the data channel for processing
func (p *Parser) parseCsv(data []byte, fname string) error {
	cif := ciface.NewParser(data)
	docs, err := cif.Parse()

	if err == nil {
		// push the docs onto the data channel
		for _, doc := range docs {
			err = p.parseFileName(fname, &doc)

			p.DataChannel <- doc
		}
	}

	return err
}

// parseFilename handles meta data extraction from filenames. It reads the pattern
// file specified with the --name-pattern argument and parses the filename according
func (p *Parser) parseFileName(fname string, doc interface{}) error {
	var err error

	if p.Cli.String("name-pattern") != "" {
		var pattern = make(map[string]interface{})

		pdoc, err := ioutil.ReadFile(p.Cli.String("name-pattern"))
		err = json.Unmarshal(pdoc, &pattern)

		if err == nil {
			if regex := regexp.MustCompile(pattern["pattern"].(string)); regex.MatchString(fname) {
				matches := regex.FindStringSubmatch(fname)
				output := pattern["output"].(string)

				// Replace the %<count> indicators in the output with the matching capture
				for i, match := range matches {
					rxp := regexp.MustCompile("(%" + strconv.Itoa(i) + ")")
					output = rxp.ReplaceAllString(output, match)
				}

				doc.(map[string]interface{})[pattern["key"].(string)] = output
			}
		}
	}

	return err
}

func (p *Parser) parsable(file string) bool {
	return !p.isConfiguration(file) && p.supportedFormat(file)
}

// @TODO test behavior for relative paths
// isConfiguration returns true if the file read matches one of the
// files specified with a configuration flag on the commandline
func (p *Parser) isConfiguration(file string) bool {
	switch file {
	case p.Cli.String("key-map"), p.Cli.String("merge"), p.Cli.String("name-pattern"), p.Cli.String("wrapper"):
		return true
	default:
		return false
	}
}

// supportedFormat returns a true if the input and format type match
func (p *Parser) supportedFormat(file string) bool {
	return p.isText(file) || p.isJson(file) || p.isCsv(file)
}

// isText returns true when a supported text flavor is provided && input format = text
func (p *Parser) isText(file string) bool {
	text := regexp.MustCompile(textFileRegex)
	return p.Cli.String("format") == "text" && text.MatchString(file)
}

// isJson returns true when a supported json flavor is detected && input format = json
func (p *Parser) isJson(file string) bool {
	json := regexp.MustCompile(jsonFileRegex)
	return p.Cli.String("format") == "json" && json.MatchString(file)
}

// isCsv returns true when a supported csv flavor is detected && input format = csv
func (p *Parser) isCsv(file string) bool {
	csv := regexp.MustCompile(csvFileRegex)
	return p.Cli.String("format") == "csv" && csv.MatchString(file)
}

func (p *Parser) replaceNewLines(data []byte, replacement string) []byte {
	newline := regexp.MustCompile(newlineRegex)
	return newline.ReplaceAll(data, []byte(replacement))
}

// handleFatal logs the error and exits the program when the error input isn't nil
func (p *Parser) handleFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
