package ghostdoc

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const (
	textFileRegex = `(?i)^.+\.txt$`
	jsonFileRegex = `(?i)^.+\.json|geojson|topojson$`
	csvFileRegex  = `(?i)^.+\.csv|tsv|txt`
)

type Parser struct {
	Cli         *cli.Context
	DataChannel *chan interface{}
	Input       []string
}

// NewParser initializes the Parser struct and passes a pointer to the caller
func NewParser(c *cli.Context, dc *chan interface{}) *Parser {
	return &Parser{
		Cli:         c,
		DataChannel: dc,
		Input:       c.Args(), // Load the commandline arguments
	}
}

// Parse the provided input
func (p *Parser) Parse() {
	// Check if the provided input is valid and trigger the file or directory parsing
	p.checkInput()
}

// checkInput performs basic argument checks and throws an error or triggers an input handler
func (p *Parser) checkInput() {
	if len(p.Input) > 0 && len(p.Input[0]) > 0 {
		p.handleInput()
	} else {
		name := p.Cli.App.Name
		log.Fatalln(name, "called without an argument. See", name, "-h for usage information.")
	}
}

// handleInput triggers file or directory parsing depending on the input type
func (p *Parser) handleInput() bool {
	for input := range p.Input {
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
			p.globDir(entry)
		} else {
			if fName := entry.Name(); p.parsable(fname) {
				p.parseFile(fName)
			}
		}
	}
}

// parseFile parses the parser input (p.Input)
func (p *Parser) parseFile(file string) bool {

}

func (p *Parser) parsable(file string) {
	return !p.isConfiguration(file) && p.supportedFormat(file)
}

// @TODO test behavior for with relative paths
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

func (p *Parser) supportedFormat(file string) bool {
	return p.isText(file) || p.isJson(file) || p.isCsv(file)
}

// isText returns true when a supported text flavor is provided && input format = text
func (p *Parser) isText(file string) bool {
	text := regexp.MustCompile(textFileRegex)
	return p.Cli.String("") == "text" && text.MatchString(file)
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

// handleFatal logs the error and exits the program when the error input isn't nil
func (p *Parser) handleFatal(err) {
	if err != nil {
		log.Fatalln(err)
	}
}
