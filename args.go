package ghostdoc

import (
	"errors"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// TypeHandler interface used to check filetypes
type TypeHandler interface {
	rawInput(string) bool
	supportedFile(string) bool
}

type ArgumentHandler struct {
	Cli     *cli.Context
	RawChan chan []byte
	RawSync *sync.WaitGroup
	TypeHandler
}

func NewArgumentHandler(c *cli.Context, raw chan []byte) *ArgumentHandler {
	return &ArgumentHandler{
		Cli:     c,
		RawChan: raw,
		RawSync: &sync.WaitGroup{},
	}
}

// hasArgs checks if any commandline arguments where provided
func (a *ArgumentHandler) hasArgs() (bool, error) {
	if len(a.Cli.Args()) == 0 {
		return false, errors.New("[Argument Error] Called without arguments: " + a.Cli.App.Name + " -h for usage info.")
	}

	return true, nil
}

// processArguments loops through all arguments and calls input handling
func (a *ArgumentHandler) processArguments() {
	for _, argument := range a.Cli.Args() {
		a.handleInput(argument)
	}
}

func (a *ArgumentHandler) handleInput(argument string) {
	if a.rawInput(argument) {
		a.RawSync.Add(1)
		a.RawChan <- []byte(argument)
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
			a.handleDiskInput(input+"/"+item.Name(), a.Cli.GlobalBool("recursive"))
		}
	} else {
		log.Println("[Argument Error]", err)
	}
}

func (a *ArgumentHandler) handleFileInput(input string) {
	if raw, err := ioutil.ReadFile(input); err == nil {
		a.RawSync.Add(1)
		go func() {
			a.RawChan <- raw
		}()
	} else {
		log.Println("[File Error]", err)
	}
}

func (a *ArgumentHandler) configuration(input string) bool {
	configuration := false

	// Grab both global and sub flags
	flags := a.Cli.GlobalFlagNames()
	flags = append(flags, a.Cli.FlagNames()...)

	// Check if the flag matches the input value
	for _, flag := range flags {
		if !configuration {
			if val := a.Cli.GlobalString(flag); val != "" {
				configuration = a.sourceCompare(input, val)
			}

			if val := a.Cli.String(flag); val != "" {
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
