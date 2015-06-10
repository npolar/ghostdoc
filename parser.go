package ghostdoc

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/npolar/ghostdoc/context"
	"github.com/npolar/ghostdoc/util"
)

// Parser typedef
type Parser struct {
	context context.GhostContext
	ParserStrategy
	argumentHandler *ArgumentHandler
	rawChan         chan *rawFile
	dataChan        chan *dataFile
}

// NewParser constructor
func NewParser(parserStrategy ParserStrategy) *Parser {
	context := parserStrategy.getContext()
	rawChan := make(chan *rawFile, context.GlobalInt("concurrency"))
	parser := &Parser{context: context, ParserStrategy: parserStrategy, rawChan: rawChan}
	parser.argumentHandler = NewArgumentHandler(parser, rawChan)
	util.ConfigureLogger(context)
	return parser
}

func (p *Parser) process() {
	var start = time.Now()
	p.dataChan = make(chan *dataFile, p.context.GlobalInt("concurrency"))
	writer := NewWriter(p.context, p.dataChan)
	writerWaitGroup, err := writer.listen()
	if err != nil {
		panic(err.Error())
	}

	p.listen()
	writerWaitGroup.Wait()

	util.SendErrorMail()
	log.Info("Stop, took: ", time.Now().Sub(start))
}

func (p *Parser) listen() {
	if ok, err := p.argumentHandler.hasArgs(); ok {
		p.argumentHandler.processArguments()
		go func() {
			for rawFile := range p.rawChan {
				p.parse(rawFile, p.dataChan)
			}
			close(p.dataChan)
		}()
	} else {
		fmt.Println(err)
	}
}
