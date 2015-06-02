package ghostdoc

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/npolar/ghostdoc/context"
)

// JSONParser typedef
type JSONParser struct {
	context     context.GhostContext
	DataChannel chan interface{}
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
func NewJSONParser(c context.GhostContext, dc chan interface{}, wg *sync.WaitGroup) *JSONParser {
	parser := &JSONParser{
		context:     c,
		DataChannel: dc,
		WaitGroup:   wg,
	}

	// Configure the argument handler and give it a channel for the raw data
	inputChan := make(chan *RawFile, c.GlobalInt("concurrency"))
	parser.ArgumentHandler = NewArgumentHandler(c, inputChan)
	// Customize the argument handler to relate to json values
	parser.TypeHandler = &JSONHandler{}

	return parser
}

func processJSON(c *cli.Context) {
	var jsonChan = make(chan interface{})
	wg := &sync.WaitGroup{}

	context := context.NewCliContext(c)

	writer := NewWriter(context, jsonChan, wg)
	if err := writer.listen(); err != nil {
		panic(err.Error())
	}

	parser := NewJSONParser(context, jsonChan, wg)
	parser.parse()

	// Wait for all go routines to finish before exiting
	wg.Wait()
}

func (jsp *JSONParser) parse() {
	if ok, err := jsp.hasArgs(); ok {
		jsp.processArguments()

		go func() {
			for {
				jsp.parseToInterface(<-jsp.RawChan)
			}
		}()

	} else {
		log.Println("[JSON] Argument Error:", err)
	}
}

// parseToInterface reads the raw json data and converts it to an interface{}
func (jsp *JSONParser) parseToInterface(raw *RawFile) {
	var jsonData interface{}

	if err := json.Unmarshal(raw.data, &jsonData); err == nil {
		if jsonData, err = jsp.parseFileName(raw.name, jsonData); err != nil {
			log.Println("[JSON] Filename parsing error!", err)
			return
		}
		jsp.WaitGroup.Add(1)
		go func(d interface{}) {
			jsp.DataChannel <- d
		}(jsonData)
	} else {
		log.Println("[JSON] Parsing error!", err)
	}
}
