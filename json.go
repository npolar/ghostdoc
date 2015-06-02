package ghostdoc

import (
	"regexp"
)

const (
	jsonFileRegex = `(?i)^.+\.json|geojson|topojson$`
)

// JSONHandler typedef
type JSONHandler struct{}

// NewJSONHandler factory
func NewJSONHandler() *JSONHandler {
	return &JSONHandler{}
}

// rawInput does a lazy check for raw inline input and returns true if matches
func (j *JSONHandler) rawInput(argument string) bool {
	rawJSON := regexp.MustCompile(`(?m)^\[|{\".+\":.+}|]$`)
	return rawJSON.MatchString(argument)
}

// supportedFile returns true if the filename meets the requirements
func (j *JSONHandler) supportedFile(filename string) bool {
	jsonFile := regexp.MustCompile(jsonFileRegex)
	return jsonFile.MatchString(filename)
}
