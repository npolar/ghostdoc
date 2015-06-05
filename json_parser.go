package ghostdoc

import (
	"encoding/json"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/npolar/ghostdoc/context"
)

// JSONParser typedef
type JSONParser struct {
	context     context.GhostContext
	DataChannel chan *dataFile
	WaitGroup   *sync.WaitGroup
	*ArgumentHandler
}

// JSONCommand cli.Command for JSON parsing
func JSONCommand() cli.Command {
	return cli.Command{
		Name:  "json",
		Usage: "parse json files",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "bulk, b",
				Usage: "enable bulk processing",
			},
			cli.StringFlag{
				Name:  "map, m",
				Usage: "rename keys. Format: '{\"oldKey\": \"newKey\"}'",
			},
		},
		Action: processJSON,
	}
}

// NewJSONParser factory
func NewJSONParser(c context.GhostContext, dc chan *dataFile, wg *sync.WaitGroup) *JSONParser {
	parser := &JSONParser{
		context:     c,
		DataChannel: dc,
		WaitGroup:   wg,
	}

	// Configure the argument handler and give it a channel for the raw data
	rawChan := make(chan *rawFile, c.GlobalInt("concurrency"))
	parser.ArgumentHandler = NewArgumentHandler(c, rawChan)
	// Customize the argument handler to relate to json values
	parser.TypeHandler = &JSONHandler{}

	return parser
}

func processJSON(c *cli.Context) {
	var dataChan = make(chan *dataFile, c.GlobalInt("concurrency"))
	wg := &sync.WaitGroup{}

	context := context.NewCliContext(c)

	writer := NewWriter(context, dataChan, wg)
	if err := writer.listen(); err != nil {
		panic(err.Error())
	}

	parser := NewJSONParser(context, dataChan, wg)
	parser.parse()

	// Wait for all go routines to finish before exiting
	wg.Wait()
}

func (jsp *JSONParser) parse() {
	if ok, err := jsp.hasArgs(); ok {
		jsp.processArguments()

		go func() {
			for rawFile := range jsp.RawChan {
				jsp.parseToInterface(rawFile)
			}
			close(jsp.DataChannel)
		}()

	} else {
		log.Error("[JSON] Argument Error:", err)
	}
}

// parseToInterface reads the raw json data and converts it to an interface{}
func (jsp *JSONParser) parseToInterface(raw *rawFile) {
	var jsonData interface{}

	if err := json.Unmarshal(raw.data, &jsonData); err == nil {
		jsp.DataChannel <- &dataFile{
			name: raw.name,
			data: jsonData.(map[string]interface{}),
		}

	} else {
		log.Error("[JSON] Parsing error!", err)
	}
}
