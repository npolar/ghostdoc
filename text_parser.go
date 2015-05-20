package ghostdoc

import (
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
)

const (
	newlineRegex = `\n|\r\n|\n\r$`
)

type TextParser struct {
	Cli         *cli.Context
	DataChannel chan interface{}
	WaitGroup   *sync.WaitGroup
	*ArgumentHandler
}

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

func NewTextParser(c *cli.Context, dc chan interface{}, wg *sync.WaitGroup) *TextParser {
	parser := &TextParser{
		Cli:         c,
		DataChannel: dc,
		WaitGroup:   wg,
	}

	// Configure the argument handler and give it a channel for the raw data
	inputChan := make(chan [][]byte, c.GlobalInt("concurrency"))
	parser.ArgumentHandler = NewArgumentHandler(c, inputChan)
	// Customize the argument handler ro relate to text values
	parser.TypeHandler = &TextHandler{}

	return parser
}

func processText(c *cli.Context) {
	var textChan = make(chan interface{})
	wg := &sync.WaitGroup{}

	// Initialize a new parser and parse the input
	parser := NewTextParser(c, textChan, wg)
	parser.parse()

	// Setup the writer for output handling
	writer := NewWriter(c, textChan, wg)
	if err := writer.Write(); err != nil {
		log.Println(err.Error())
	}
	// Wait for all go routines to finish before exiting
	wg.Wait()
}

func (tp *TextParser) parse() {
	if ok, err := tp.hasArgs(); ok {
		tp.processArguments()

		go func() {
			for {
				tp.parseToInterface(<-tp.RawChan)
				tp.RawSync.Done()
			}
		}()

		tp.RawSync.Wait()
	} else {
		log.Println("[Text] Argument error:", err)
	}
}

func (tp *TextParser) parseToInterface(raw [][]byte) {
	var dataMap = make(map[string]interface{})
	var textIface interface{}

	text := tp.replaceNewLines(raw[1], " ")
	dataMap[tp.Cli.String("key")] = strings.TrimSpace(string(text))
	textIface = dataMap

	//	textIface, err := tp.parseFileName(fname, textIface)

	tp.WaitGroup.Add(1)
	go func(d interface{}) {
		tp.DataChannel <- d
	}(textIface)
}

func (tp *TextParser) replaceNewLines(data []byte, replacement string) []byte {
	newline := regexp.MustCompile(newlineRegex)
	return newline.ReplaceAll(data, []byte(replacement))
}
