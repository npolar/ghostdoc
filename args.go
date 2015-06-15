package ghostdoc

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/npolar/ghostdoc/context"
)

type rawFile struct {
	name string
	data []byte
}

// ArgumentHandler typdef
type ArgumentHandler struct {
	context context.GhostContext
	rawChan chan *rawFile
	parser  *Parser
}

// NewArgumentHandler factory
func NewArgumentHandler(parser *Parser, rawChan chan *rawFile) *ArgumentHandler {
	return &ArgumentHandler{
		context: parser.getContext(),
		rawChan: rawChan,
		parser:  parser,
	}
}

// HasArgs checks if any commandline arguments where provided
func (a *ArgumentHandler) hasArgs() (bool, error) {
	if len(a.context.Args()) == 0 && !a.hasPipe() {
		return false, errors.New("[Argument Error] Called without arguments: " + a.context.Cli().App.Name + " -h for usage info.")
	}

	return true, nil
}

// ProcessArguments loops through all arguments and calls input handling
func (a *ArgumentHandler) processArguments() {
	go func() {
		if a.hasPipe() {
			bytes, _ := ioutil.ReadAll(os.Stdin)
			a.handleInput(string(bytes))
		} else {
			for _, argument := range a.context.Args() {
				a.handleInput(argument)
			}
		}
		close(a.rawChan)
	}()
}

func (a *ArgumentHandler) hasPipe() bool {
	fi, err := os.Stdin.Stat()
	return !(fi.Mode()&os.ModeNamedPipe == 0) && err == nil
}

func (a *ArgumentHandler) handleInput(argument string) {
	if a.parser.isRawInput(argument) {
		log.Info("Parsing raw input")
		data := &rawFile{
			name: a.context.GlobalString("filename"),
			data: []byte(argument),
		}
		a.rawChan <- data
	} else {
		a.handleDiskInput(argument, false)
	}
}

func (a *ArgumentHandler) handleDiskInput(argument string, recursive bool) {
	if state, err := os.Stat(argument); err == nil {
		if state.IsDir() {
			a.globDir(argument)
		} else if !a.configuration(argument) && a.parser.isSupportedFile(argument) {
			a.handleFileInput(argument)
		} else {
			log.Warn("[Unsupported Filetype] Skipping:", argument)
		}
	} else {
		log.WithFields(log.Fields{"input": argument}).Warn("[Input Error] ", err)
	}
}

func (a *ArgumentHandler) globDir(input string) {
	if dirList, err := ioutil.ReadDir(input); err == nil {
		for _, item := range dirList {
			a.handleDiskInput(input+"/"+item.Name(), a.context.GlobalBool("recursive"))
		}
	} else {
		log.WithFields(log.Fields{"input": input}).Error("[Argument Error] ", err)
	}
}

func (a *ArgumentHandler) handleFileInput(input string) {
	if raw, err := ioutil.ReadFile(input); err == nil {
		log.Info("Parsing ", input)
		data := &rawFile{
			name: input,
			data: raw,
		}
		a.rawChan <- data
	} else {
		log.WithFields(log.Fields{"input": input}).Error("[File Error] ", err)
	}
}

func (a *ArgumentHandler) configuration(input string) bool {
	configuration := false

	// Grab both global and sub flags
	flags := a.context.Cli().GlobalFlagNames()
	flags = append(flags, a.context.FlagNames()...)

	// Check if the flag matches the input value
	for _, flag := range flags {
		if !configuration {
			if val := a.context.String(flag); val != "" {
				configuration = a.sourceCompare(input, val)
			}

			if val := a.context.GlobalString(flag); val != "" {
				configuration = a.sourceCompare(input, val)
			}
		}
	}

	return configuration
}

func (a *ArgumentHandler) sourceCompare(firstPath string, secondPath string) bool {
	firstPath, _ = filepath.Abs(firstPath)
	secondPath, _ = filepath.Abs(secondPath)
	return firstPath == secondPath
}
