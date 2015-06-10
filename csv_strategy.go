package ghostdoc

import (
	"io/ioutil"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/npolar/ciface"
	"github.com/npolar/ghostdoc/context"
	"github.com/npolar/ghostdoc/util"
)

// CsvStrategy typedef
type CsvStrategy struct {
	context   context.GhostContext
	delimiter string
}

const (
	csvExtRegex = `(?i)^.+\.csv|tsv|txt`
)

// NewCsvStrategy factory
func NewCsvStrategy(context context.GhostContext) *CsvStrategy {
	return &CsvStrategy{context: context, delimiter: context.String("delimiter")}
}

// rawInput does a lazy check for raw inline input and returns true if matches
func (c *CsvStrategy) isRawInput(argument string) bool {
	rawCsv := regexp.MustCompile(`(?m)^(.*)` + c.delimiter + `(.*)$`)
	return rawCsv.MatchString(argument)
}

// supportedFile returns true if the filename meets the requirements
func (c *CsvStrategy) isSupportedFile(filename string) bool {
	csvFile := regexp.MustCompile(csvExtRegex)
	return csvFile.MatchString(filename)
}

func (c *CsvStrategy) getContext() context.GhostContext {
	return c.context
}

func (c *CsvStrategy) parse(rawFile *rawFile, dataChan chan *dataFile) {
	cif := ciface.NewParser(rawFile.data)
	cif.Skip = c.context.Int("skip")

	if header := c.context.String("header"); header != "" {
		hfile, err := ioutil.ReadFile(header)
		if err != nil {
			cif.Header = util.StringToSlice(header)
		} else {
			cif.Header = util.StringToSlice(string(hfile))
		}
	}

	delimiterRune, _, _, _ := strconv.UnquoteChar(c.context.String("delimiter"), '"')
	cif.Reader.Comma = delimiterRune

	commentRune, _, _, _ := strconv.UnquoteChar(c.context.String("comment"), '"')
	cif.Reader.Comment = commentRune

	docs, err := cif.Parse()

	// push the docs onto the data channel
	for _, doc := range docs {
		dataChan <- &dataFile{
			name: rawFile.name,
			data: doc.(map[string]interface{}),
		}
	}

	if err != nil {
		log.Error("[Parsing error]", err)
	}
}
