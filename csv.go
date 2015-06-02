package ghostdoc

import "regexp"

const (
	csvExtRegex = `(?i)^.+\.csv|tsv|txt`
)

// CsvHandler typedef
type CsvHandler struct {
	Delimiter string
}

// NewCsvHandler factory
func NewCsvHandler(delimiter string) *CsvHandler {
	return &CsvHandler{
		Delimiter: delimiter,
	}
}

// rawInput does a lazy check for raw inline input and returns true if matches
func (c *CsvHandler) rawInput(argument string) bool {
	rawCsv := regexp.MustCompile(`(?m)^(.*)` + c.Delimiter + `(.*)$`)
	return rawCsv.MatchString(argument)
}

// supportedFile returns true if the filename meets the requirements
func (c *CsvHandler) supportedFile(filename string) bool {
	csvFile := regexp.MustCompile(csvExtRegex)
	return csvFile.MatchString(filename)
}
