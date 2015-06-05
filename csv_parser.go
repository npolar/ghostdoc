package ghostdoc

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/npolar/ciface"
	"github.com/npolar/ghostdoc/context"
	"github.com/npolar/ghostdoc/util"
)

// CsvParser typedef
type CsvParser struct {
	context     context.GhostContext
	DataChannel chan *dataFile
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

func newCsvParser(c context.GhostContext, dc chan *dataFile, wg *sync.WaitGroup) *CsvParser {
	parser := &CsvParser{
		context:     c,
		DataChannel: dc,
		WaitGroup:   wg,
	}

	// Configure the argument handler and give it a channel for the raw data
	rawChan := make(chan *rawFile, c.GlobalInt("concurrency"))
	parser.ArgumentHandler = NewArgumentHandler(c, rawChan)
	// Customize the argument handler to relate to csv values
	parser.TypeHandler = &CsvHandler{
		Delimiter: c.String("delimiter"),
	}

	return parser
}

func processCsv(c *cli.Context) {
	var dataChan = make(chan *dataFile, c.GlobalInt("concurrency"))
	wg := &sync.WaitGroup{}

	context := context.NewCliContext(c)

	writer := NewWriter(context, dataChan, wg)
	if err := writer.listen(); err != nil {
		panic(err.Error())
	}

	csvParser := newCsvParser(context, dataChan, wg)
	csvParser.parse()

	// Wait for all go routines to finish before exiting
	wg.Wait()
}

func (csv *CsvParser) parse() {
	if ok, err := csv.hasArgs(); ok {
		csv.processArguments()

		go func() {
			for rawFile := range csv.RawChan {
				csv.parseToInterface(rawFile)
			}
			close(csv.DataChannel)
		}()
	} else {
		fmt.Println(err)
	}
}

func (csv *CsvParser) parseToInterface(raw *rawFile) {
	cif := ciface.NewParser(raw.data)
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
		csv.DataChannel <- &dataFile{
			name: raw.name,
			data: doc.(map[string]interface{}),
		}
	}

	if err != nil {
		log.Error("[Parsing error]", err)
	}
}
