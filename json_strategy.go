package ghostdoc

import (
	"encoding/json"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/npolar/ghostdoc/context"
)

const (
	jsonFileRegex = `(?i)^.+\.json|geojson|topojson$`
)

// JSONStrategy typedef
type JSONStrategy struct {
	context   context.GhostContext
	delimiter string
}

// NewJSONStrategy factory
func NewJSONStrategy(context context.GhostContext) *JSONStrategy {
	return &JSONStrategy{context: context, delimiter: context.String("delimiter")}
}

// rawInput does a lazy check for raw inline input and returns true if matches
func (j *JSONStrategy) isRawInput(argument string) bool {
	rawJSON := regexp.MustCompile(`(?m)^\[|{\".+\":.+}|]$`)
	return rawJSON.MatchString(argument)
}

// supportedFile returns true if the filename meets the requirements
func (j *JSONStrategy) isSupportedFile(filename string) bool {
	jsonFile := regexp.MustCompile(jsonFileRegex)
	return jsonFile.MatchString(filename)
}

func (j *JSONStrategy) getContext() context.GhostContext {
	return j.context
}

func (j *JSONStrategy) parse(rawFile *rawFile, dataChan chan *dataFile) {
	var jsonData interface{}

	if err := json.Unmarshal(rawFile.data, &jsonData); err == nil {
		dataChan <- &dataFile{
			name: rawFile.name,
			data: jsonData.(map[string]interface{}),
		}

	} else {
		log.Error("[JSON] Parsing error!", err)
	}
}
