package ghostdoc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/npolar/ghostdoc/context"
)

// ArgumentHandler typdef
type ArgumentHandler struct {
	context context.GhostContext
	RawChan chan [][]byte
	RawSync *sync.WaitGroup
	TypeHandler
}

// NewArgumentHandler factory
func NewArgumentHandler(c context.GhostContext, raw chan [][]byte) *ArgumentHandler {
	return &ArgumentHandler{
		context: c,
		RawChan: raw,
		RawSync: &sync.WaitGroup{},
	}
}

// HasArgs checks if any commandline arguments where provided
func (a *ArgumentHandler) hasArgs() (bool, error) {
	if len(a.context.Args()) == 0 && !a.hasPipe() {
		return false, errors.New("[Argument Error] Called without arguments: " + a.context.Cli().App.Name + " -h for usage info.")
	}

	return true, nil
}

// ParseFileName handles meta data extraction from filenames. It reads the pattern
// file specified with the --name-pattern argument and parses the filename according
func (a *ArgumentHandler) parseFileName(fname string, doc interface{}) (interface{}, error) {
	var err error
	if pat := a.context.GlobalString("name-pattern"); pat != "" {
		var pattern = make(map[string]interface{})
		var pdoc []byte
		if pdoc, err = ioutil.ReadFile(pat); err != nil {
			pdoc = []byte(pat)
		}

		if err = json.Unmarshal(pdoc, &pattern); err == nil {
			if pRgx := regexp.MustCompile(pattern["pattern"].(string)); pRgx.MatchString(fname) {
				matches := pRgx.FindStringSubmatch(fname)
				outputB, _ := json.Marshal(pattern["output"])
				output := string(outputB)

				// Replace the %<count> indicators in the output with the matching capture
				for i, match := range matches {
					rxp := regexp.MustCompile("(%" + strconv.Itoa(i) + ")")
					output = rxp.ReplaceAllString(output, match)
				}

				var jsonData map[string]interface{}
				if err := json.Unmarshal([]byte(output), &jsonData); err == nil {
					for key, val := range jsonData {
						doc.(map[string]interface{})[key] = val
					}
				}
			}
		}
	}

	if err != nil {
		err = errors.New("name-pattern: " + err.Error())
	}

	return doc, err
}

// ProcessArguments loops through all arguments and calls input handling
func (a *ArgumentHandler) processArguments() {
	if a.context.GlobalBool("quiet") {
		log.SetOutput(ioutil.Discard)
	}

	if a.hasPipe() {
		bytes, _ := ioutil.ReadAll(os.Stdin)
		a.handleInput(string(bytes))
	} else {
		for _, argument := range a.context.Args() {
			a.handleInput(argument)
		}
	}
}

func (a *ArgumentHandler) hasPipe() bool {
	fi, err := os.Stdin.Stat()
	return !(fi.Mode()&os.ModeNamedPipe == 0) && err == nil
}

func (a *ArgumentHandler) handleInput(argument string) {
	if a.rawInput(argument) {
		data := make([][]byte, 2)
		data[0] = []byte(a.context.GlobalString("filename"))
		data[1] = []byte(argument)
		a.RawSync.Add(1)
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
		data := make([][]byte, 2)
		data[0] = []byte(input)
		data[1] = raw
		a.RawSync.Add(1)
		go func() {
			a.RawChan <- data
		}()
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
