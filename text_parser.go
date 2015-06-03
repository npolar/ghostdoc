package ghostdoc

import (
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/npolar/ghostdoc/context"
)

const (
	newlineRegex = `\n|\r\n|\n\r$`
)

// TextParser typedef
type TextParser struct {
	context     context.GhostContext
	DataChannel chan *dataFile
	WaitGroup   *sync.WaitGroup
	*ArgumentHandler
}

// TextCommand cli.Command for parsing text
func TextCommand() cli.Command {
	return cli.Command{
		Name:    "text",
		Aliases: []string{"txt"},
		Usage:   "parse text data",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "key, k",
				Usage: "the key used to identify the text segment",
			},
			cli.StringFlag{
				Name:  "pattern, p",
				Usage: "use a pattern file to specify which text segments should be extracted",
			},
		},
		Action: processText,
	}
}

// NewTextParser factory
func NewTextParser(c context.GhostContext, dc chan *dataFile, wg *sync.WaitGroup) *TextParser {
	parser := &TextParser{
		context:     c,
		DataChannel: dc,
		WaitGroup:   wg,
	}

	// Configure the argument handler and give it a channel for the raw data
	inputChan := make(chan *rawFile, c.GlobalInt("concurrency"))
	parser.ArgumentHandler = NewArgumentHandler(c, inputChan)
	// Customize the argument handler ro relate to text values
	parser.TypeHandler = &TextHandler{}

	return parser
}

func processText(c *cli.Context) {
	var dataChan = make(chan *dataFile, c.GlobalInt("concurrency"))
	wg := &sync.WaitGroup{}

	context := context.NewCliContext(c)

	// Setup the writer for output handling
	writer := NewWriter(context, dataChan, wg)
	if err := writer.listen(); err != nil {
		panic(err.Error())
	}

	// Initialize a new parser and parse the input
	parser := NewTextParser(context, dataChan, wg)
	parser.parse()

	// Wait for all go routines to finish before exiting
	wg.Wait()
}

func (tp *TextParser) parse() {
	if ok, err := tp.hasArgs(); ok {
		tp.processArguments()

		go func() {
			for rawFile := range tp.RawChan {
				tp.parseToInterface(rawFile)
			}
			close(tp.DataChannel)
		}()

	} else {
		log.Println("[Text] Argument error:", err)
	}
}

func (tp *TextParser) parseToInterface(raw *rawFile) {
	var dataMap = make(map[string]interface{})

	text := tp.replaceNewLines(raw.data, " ")
	dataMap[tp.context.String("key")] = strings.TrimSpace(string(text))

	tp.DataChannel <- &dataFile{
		name: raw.name,
		data: dataMap,
	}
}

func (tp *TextParser) replaceNewLines(data []byte, replacement string) []byte {
	newline := regexp.MustCompile(newlineRegex)
	return newline.ReplaceAll(data, []byte(replacement))
}
