package ghostdoc

import (
	"regexp"
)

const (
	jsonFileRegex = `(?i)^.+\.json|geojson|topojson$`
)

type JsonHandler struct{}

func NewJsonHandler() *JsonHandler {
	return &JsonHandler{}
}

// rawInput does a lazy check for raw inline input and returns true if matches
func (j *JsonHandler) rawInput(argument string) bool {
	rawJson := regexp.MustCompile(`(?m)^\[|{\".+\":.+}|]$`)
	return rawJson.MatchString(argument)
}

// supportedFile returns true if the filename meets the requirements
func (j *JsonHandler) supportedFile(filename string) bool {
	jsonFile := regexp.MustCompile(jsonFileRegex)
	return jsonFile.MatchString(filename)
}
