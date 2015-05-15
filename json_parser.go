package ghostdoc

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"log"
	"sync"
)

type JsonParser struct {
	Cli         *cli.Context
	DataChannel chan interface{}
	WaitGroup   *sync.WaitGroup
	*ArgumentHandler
}

func JsonCommand() cli.Command {
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
		Action: processJson,
	}
}

func NewJsonParser(c *cli.Context, dc chan interface{}, wg *sync.WaitGroup) *JsonParser {
	parser := &JsonParser{
		Cli:         c,
		DataChannel: dc,
		WaitGroup:   wg,
	}

	// Configure the argument handler and give it a channel for the raw data
	inputChan := make(chan []byte, c.GlobalInt("concurrency"))
	parser.ArgumentHandler = NewArgumentHandler(c, inputChan)
	// Customize the argument handler to relate to json values
	parser.TypeHandler = &JsonHandler{}

	return parser
}

func processJson(c *cli.Context) {
	var jsonChan = make(chan interface{})
	wg := &sync.WaitGroup{}

	parser := NewJsonParser(c, jsonChan, wg)
	parser.parse()

	writer := NewWriter(c, jsonChan, wg)
	if err := writer.Write(); err != nil {
		log.Println(err.Error())
	}
	// Wait for all go routines to finish before exiting
	wg.Wait()
}

func (jsp *JsonParser) parse() {
	if ok, err := jsp.hasArgs(); ok {
		jsp.processArguments()

		go func() {
			for {
				jsp.parseToInterface(<-jsp.RawChan)
				jsp.RawSync.Done()
			}
		}()

		jsp.RawSync.Wait()
	} else {
		log.Println(err)
	}
}

// parseToInterface reads the raw json data and converts it to an interface{}
func (jsp *JsonParser) parseToInterface(raw []byte) {
	var jsonData interface{}
	if err := json.Unmarshal(raw, &jsonData); err == nil {
		jsp.WaitGroup.Add(1)
		go func(d interface{}) {
			jsp.DataChannel <- d
		}(jsonData)
	} else {
		log.Println("[JSON] Parsing error!", err)
	}
}
