package ghostdoc

import (
	"regexp"
	"strings"

	"github.com/npolar/ghostdoc/context"
)

// TextStrategy typedef
type TextStrategy struct {
	context context.GhostContext
}

const (
	textFileRegex = `(?i)^.+\.txt$`
	newlineRegex  = `\n`
)

// NewTextStrategy factory
func NewTextStrategy(context context.GhostContext) *TextStrategy {
	return &TextStrategy{context: context}
}

// rawInput does a lazy check for raw inline input and returns true if matches
func (t *TextStrategy) isRawInput(argument string) bool {
	rawText := regexp.MustCompile(`(?m)^\[|{\".+\":.+}|]$`)
	return rawText.MatchString(argument)
}

// supportedFile returns true if the filename meets the requirements
func (t *TextStrategy) isSupportedFile(filename string) bool {
	textFile := regexp.MustCompile(textFileRegex)
	return textFile.MatchString(filename)
}

func (t *TextStrategy) getContext() context.GhostContext {
	return t.context
}

func (t *TextStrategy) parse(rawFile *rawFile, dataChan chan *dataFile) {
	var dataMap = make(map[string]interface{})
	text := t.replaceNewLines(rawFile.data, " ")
	dataMap[t.context.String("key")] = strings.TrimSpace(string(text))

	dataChan <- &dataFile{
		name: rawFile.name,
		data: dataMap,
	}
}

func (t *TextStrategy) replaceNewLines(data []byte, replacement string) []byte {
	newline := regexp.MustCompile(newlineRegex)
	return newline.ReplaceAll(data, []byte(replacement))
}
