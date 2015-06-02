package ghostdoc

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/npolar/ciface"
	"github.com/npolar/ghostdoc/context"
	"github.com/npolar/ghostdoc/util"
)

// CsvParser typedef
type CsvParser struct {
	context     context.GhostContext
	DataChannel chan interface{}
	WaitGroup   *sync.WaitGroup
	*ArgumentHandler
}

// CsvCommand specifies the cli interface for the csv parser
func CsvCommand() cli.Command {
	return cli.Command{
		Name:  "csv",
		Usage: "Parse delimiter separated value files (csv, tsv, etc...)",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "delimiter, d",
				Value: ",",
				Usage: "Set the delimiter char.",
			},
			cli.StringFlag{
				Name:  "comment, c",
				Value: "#",
				Usage: "Set the comment char.",
			},
			cli.StringFlag{
				Name:  "header, hd",
				Usage: "Configure data header. If not set the first data line will be used.",
			},
			cli.IntFlag{
				Name:  "skip, s",
				Usage: "Specify the number of lines to skip before parsing. [NOTE] Blank lines are ignored by the parser and should not be skipped.",
			},
		},
		Action: processCsv,
	}
}

func newCsvParser(c context.GhostContext, dc chan interface{}, wg *sync.WaitGroup) *CsvParser {
	parser := &CsvParser{
		context:     c,
		DataChannel: dc,
		WaitGroup:   wg,
	}

	// Configure the argument handler and give it a channel for the raw data
	inputChan := make(chan [][]byte, c.GlobalInt("concurrency"))
	parser.ArgumentHandler = NewArgumentHandler(c, inputChan)
	// Customize the argument handler to relate to csv values
	parser.TypeHandler = &CsvHandler{
		Delimiter: c.String("delimiter"),
	}

	return parser
}

func processCsv(c *cli.Context) {
	var csvChan = make(chan interface{})
	wg := &sync.WaitGroup{}

	context := context.NewCliContext(c)

	csvParser := newCsvParser(context, csvChan, wg)
	csvParser.parse()

	writer := NewWriter(context, csvChan, wg)
	if err := writer.Write(); err != nil {
		log.Println(err.Error())
	}
	// Wait for all go routines to finish before exiting
	wg.Wait()
}

func (csv *CsvParser) parse() {
	if ok, err := csv.hasArgs(); ok {
		csv.processArguments()

		go func() {
			for {
				csv.parseToInterface(<-csv.RawChan)
				csv.RawSync.Done()
			}
		}()

		csv.RawSync.Wait()

	} else {
		fmt.Println(err)
	}
}

func (csv *CsvParser) parseToInterface(raw [][]byte) {
	cif := ciface.NewParser(raw[1])
	cif.Skip = csv.context.Int("skip")

	if header := csv.context.String("header"); header != "" {
		hfile, err := ioutil.ReadFile(header)
		if err != nil {
			cif.Header = util.StringToSlice(header)
		} else {
			cif.Header = util.StringToSlice(string(hfile))
		}
	}

	delimiterRune, _, _, _ := strconv.UnquoteChar(csv.context.String("delimiter"), '"')
	cif.Reader.Comma = delimiterRune

	commentRune, _, _, _ := strconv.UnquoteChar(csv.context.String("comment"), '"')
	cif.Reader.Comment = commentRune

	docs, err := cif.Parse()

	// push the docs onto the data channel
	for _, doc := range docs {
		doc, err = csv.parseFileName(string(raw[0]), doc)
		csv.WaitGroup.Add(1)

		go func(d interface{}) {
			csv.DataChannel <- d
		}(doc)

	}

	if err != nil {
		log.Println("[Parsing error]", err)
	}
}
