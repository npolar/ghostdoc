package ghostdoc

import (
	"regexp"
)

const (
	textFileRegex = `(?i)^.+\.txt$`
)

type TextHandler struct{}

func NewTextHandler() *TextHandler {
	return &TextHandler{}
}

// rawInput does a lazy check for raw inline input and returns true if matches
func (j *TextHandler) rawInput(argument string) bool {
	rawText := regexp.MustCompile(`^\[|{\".+\":.+}|]$`)
	return rawText.MatchString(argument)
}

// supportedFile returns true if the filename meets the requirements
func (j *TextHandler) supportedFile(filename string) bool {
	textFile := regexp.MustCompile(textFileRegex)
	return textFile.MatchString(filename)
}
