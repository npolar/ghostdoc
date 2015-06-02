package ghostdoc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"code.google.com/p/go-uuid/uuid"
	"github.com/npolar/ghostdoc/context"
	"github.com/npolar/ghostdoc/util"
)

const (
	jsonRegex = `^\{\".+\"\:.+\}$`
)

type mapper func(data map[string]interface{}) (map[string]interface{}, error)

// Writer type definition
type Writer struct {
	context     context.GhostContext
	DataChannel chan interface{}
	WaitGroup   *sync.WaitGroup
	Js          *Js
	Validator   *Validator
}

// NewWriter initialises a new Writer and return a pointer to it
func NewWriter(c context.GhostContext, dc chan interface{}, wg *sync.WaitGroup) *Writer {
	return &Writer{
		context:     c,
		DataChannel: dc,
		WaitGroup:   wg,
		Js:          NewJs(c),
		Validator:   NewValidator(c),
	}
}

// Write applies the configured output modifiers and then writes
// the result to the configured output channel [stdout|files|http]
func (w *Writer) Write() error {
	err := w.createOutputDir()

	go func() {
		for {
			data := <-w.DataChannel
			dataMap := data.(map[string]interface{})
			dataMap, err = w.applyMappers(dataMap)

			if err == nil {
				err = w.Validator.validate(dataMap)
			}

			if err == nil {
				err = w.publishData(dataMap)
			} else {
				log.Println(err.Error())
			}

			w.WaitGroup.Done()
		}
	}()

	return err
}

func (w *Writer) applyMappers(dataMap map[string]interface{}) (map[string]interface{}, error) {
	var err error
	mappers := []mapper{
		w.includeKeys,
		w.excludeKeys,
		w.mapKeys,
		w.mergeData,
		w.wrapData,
		w.injectUUID,
		w.Js.runJs}

	for _, fn := range mappers {
		dataMap, err = fn(dataMap)
		if err != nil {
			err = errors.New("[Writer error] " + err.Error())
			break
		}
	}
	return dataMap, err
}

func (w *Writer) includeKeys(data map[string]interface{}) (map[string]interface{}, error) {
	var err error
	includeData := data

	if includes := w.context.GlobalString("include"); includes != "" {
		includesSlice := util.StringToSlice(includes)
		includeData = make(map[string]interface{})
		for _, key := range includesSlice {
			includeData[key] = data[key]
		}
	}

	return includeData, err
}

func (w *Writer) excludeKeys(data map[string]interface{}) (map[string]interface{}, error) {
	var err error

	if excludes := w.context.GlobalString("exclude"); excludes != "" {
		excludesSlice := util.StringToSlice(excludes)
		for _, key := range excludesSlice {
			delete(data, key)
		}
	}

	return data, err
}

func (w *Writer) mapKeys(data map[string]interface{}) (map[string]interface{}, error) {
	var err error
	dataMap := data

	if keyMap := w.context.GlobalString("key-map"); keyMap != "" {
		if mapping, mapErr := w.readData(keyMap); mapErr == nil {
			for key, val := range mapping {
				dataMap[val.(string)] = dataMap[key]
				delete(dataMap, key)
			}
		} else {
			err = mapErr
		}
	}

	return dataMap, err
}

func (w *Writer) wrapData(data map[string]interface{}) (map[string]interface{}, error) {
	var err error

	if wrap := w.context.GlobalString("wrapper"); wrap != "" {
		if wrapper, dataErr := w.readData(wrap); dataErr == nil {
			key := w.context.GlobalString("payload-key")
			wrapper[key] = data
			data = wrapper
		} else {
			err = dataErr
		}
	}

	return data, err
}

func (w *Writer) mergeData(data map[string]interface{}) (map[string]interface{}, error) {
	var err error

	if merge := w.context.GlobalString("merge"); merge != "" {
		if padding, dataError := w.readData(merge); dataError == nil {

			for key, val := range padding {
				data[key] = val
			}
		} else {
			err = dataError
		}
	}

	return data, err
}

func (w *Writer) injectUUID(data map[string]interface{}) (map[string]interface{}, error) {
	var err error
	keys := w.context.GlobalString("uuid-keys")

	if w.context.GlobalBool("uuid") || keys != "" {
		idData := data
		if keys != "" {
			idData = make(map[string]interface{})
			keysSlice := util.StringToSlice(keys)
			for _, key := range keysSlice {
				var ok bool
				if idData[key], ok = data[key]; !ok {
					return data, errors.New("Could not build UUID on key: " + key)
				}
			}
		}
		if doc, jsonError := json.Marshal(idData); jsonError == nil {
			data["id"] = w.generateUUID(doc)
		} else {
			err = jsonError
		}
	}

	return data, err
}

func (w *Writer) generateUUID(input []byte) string {
	id := uuid.NewSHA1(uuid.NameSpace_DNS, input)
	return id.String()
}

// createOutputDir checks the output flag and creates the
// directory (filemode 0666) on the filesystem if not present
func (w *Writer) createOutputDir() error {
	var err error

	if output := w.context.GlobalString("output"); output != "" {
		if state, statErr := os.Stat(output); state == nil {
			err = os.Mkdir(output, 0755)
		} else {
			err = statErr
		}
	}

	return err
}

// publishData grabs/generates the id of the data and converts it to a json
// document. Afterwards it calls writeFile and httpRequest methods
func (w *Writer) publishData(data map[string]interface{}) error {
	var err error
	var id string

	if doc, jsonErr := json.MarshalIndent(data, "", "  "); jsonErr == nil {
		if data["id"] != nil {
			id = data["id"].(string)
		} else {
			id = w.generateUUID(doc)
		}

		log.Println(id, string(doc))

		err = w.writeFile(doc, id)
		err = w.httpRequest(doc, id)
	} else {
		err = jsonErr
	}

	return err
}

// writeFile dumps the documents as files in the specified output dir
func (w *Writer) writeFile(doc []byte, id string) error {
	var err error

	if output := w.context.GlobalString("output"); output != "" {
		if path, pathErr := filepath.Abs(output); pathErr == nil {
			err = ioutil.WriteFile(path+"/"+id+".json", doc, 0755)
		} else {
			err = pathErr
		}
	}

	return err
}

// httpRequest performs the request defined in http-verb against
// the configured address. Default operation is POST
func (w *Writer) httpRequest(doc []byte, id string) error {
	var err error

	if addr := w.context.GlobalString("address"); addr != "" {
		if uri, uriErr := url.Parse(addr); uriErr == nil {
			client := &http.Client{}

			byteReader := bytes.NewReader(doc)

			if req, httpErr := http.NewRequest(w.context.GlobalString("http-verb"), uri.String(), byteReader); httpErr == nil {
				req.Header.Set("Content-Type", "application/json")
				var resp *http.Response
				if resp, err = client.Do(req); err == nil {
					log.Println("HTTP", w.context.GlobalString("http-verb"), "Response:", resp.Status)
				}

			} else {
				err = httpErr
			}
		} else {
			err = uriErr
		}
	}
	return err
}

func (w *Writer) readData(input string) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	var err error

	if w.jsonInput(input) {
		err = json.Unmarshal([]byte(input), &data)
	} else {
		err = w.parseJSONFile(input, &data)
	}

	return data, err
}

func (w *Writer) parseJSONFile(file string, data *map[string]interface{}) error {
	raw, err := ioutil.ReadFile(file)

	if err == nil {
		err = json.Unmarshal(raw, &data)
	}

	return err
}

func (w *Writer) jsonInput(input string) bool {
	jrxp := regexp.MustCompile(jsonRegex)
	return jrxp.MatchString(input)
}
