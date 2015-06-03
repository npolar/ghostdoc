package ghostdoc

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/npolar/ghostdoc/context"
)

type rawFile struct {
	name string
	data []byte
}

// ArgumentHandler typdef
type ArgumentHandler struct {
	context context.GhostContext
	RawChan chan *rawFile
	TypeHandler
}

// NewArgumentHandler factory
func NewArgumentHandler(c context.GhostContext, raw chan *rawFile) *ArgumentHandler {
	return &ArgumentHandler{
		context: c,
		RawChan: raw,
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
	if a.context.GlobalBool("quiet") {
		log.SetOutput(ioutil.Discard)
	}
	go func() {
		if a.hasPipe() {
			bytes, _ := ioutil.ReadAll(os.Stdin)
			a.handleInput(string(bytes))
		} else {
			for _, argument := range a.context.Args() {
				a.handleInput(argument)
			}
		}
		close(a.RawChan)
	}()
}

func (a *ArgumentHandler) hasPipe() bool {
	fi, err := os.Stdin.Stat()
	return !(fi.Mode()&os.ModeNamedPipe == 0) && err == nil
}

func (a *ArgumentHandler) handleInput(argument string) {
	if a.rawInput(argument) {
		data := &rawFile{
			name: a.context.GlobalString("filename"),
			data: []byte(argument),
		}
		a.RawChan <- data
	} else {
		a.handleDiskInput(argument, false)
	}
}

func (a *ArgumentHandler) handleDiskInput(argument string, recursive bool) {
	if state, err := os.Stat(argument); err == nil {
		if state.IsDir() {
			a.globDir(argument)
		} else if !a.configuration(argument) && a.supportedFile(argument) {
			a.handleFileInput(argument)
		} else {
			log.Println("[Unsupported Filetype] Skipping:", argument)
		}
	} else {
		log.Println("[Input Error]", err)
	}
}

func (a *ArgumentHandler) globDir(input string) {
	if dirList, err := ioutil.ReadDir(input); err == nil {
		for _, item := range dirList {
			a.handleDiskInput(input+"/"+item.Name(), a.context.GlobalBool("recursive"))
		}
	} else {
		log.Println("[Argument Error]", err)
	}
}

func (a *ArgumentHandler) handleFileInput(input string) {
	if raw, err := ioutil.ReadFile(input); err == nil {
		data := &rawFile{
			name: input,
			data: raw,
		}
		a.RawChan <- data
	} else {
		log.Println("[File Error]", err)
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
