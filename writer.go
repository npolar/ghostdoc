package ghostdoc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/cli"
)

const (
	jsonRegex = `^\{\".+\"\:.+\}$`
)

type Writer struct {
	Cli         *cli.Context
	DataChannel chan interface{}
	WaitGroup   *sync.WaitGroup
}

func NewWriter(c *cli.Context, dc chan interface{}, wg *sync.WaitGroup) *Writer {
	return &Writer{
		Cli:         c,
		DataChannel: dc,
		WaitGroup:   wg,
	}
}

// Write applies the configured output modifiers and then writes
// the result to the configured output channel [stdout|files|http]
func (w *Writer) Write() error {
	err := w.createOutputDir()

	go func() {
		for {
			data := <-w.DataChannel

			data, err = w.excludeKeys(data)
			data, err = w.mapKeys(data)
			data, err = w.mergeData(data)
			data, err = w.wrapData(data)
			data, err = w.injectUuid(data)

			if err == nil {
				err = w.publishData(data)
			}

			w.WaitGroup.Done()
		}
	}()

	return err
}

func (w *Writer) excludeKeys(data interface{}) (interface{}, error) {
	var err error

	if excludes := w.Cli.String("exclude"); excludes != "" {
		excludesSlice := strings.Split(excludes, ",")
		for _, key := range excludesSlice {
			delete(data.(map[string]interface{}), key)
		}
	}

	return data, err
}

func (w *Writer) mapKeys(data interface{}) (interface{}, error) {
	var err error
	dataMap := data.(map[string]interface{})

	if keyMap := w.Cli.String("key-map"); keyMap != "" {
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

func (w *Writer) wrapData(data interface{}) (interface{}, error) {
	var err error

	if wrap := w.Cli.String("wrapper"); wrap != "" {
		if wrapper, dataErr := w.readData(wrap); dataErr == nil {
			key := w.Cli.String("payload-key")
			wrapper[key] = data
			data = wrapper
		} else {
			err = dataErr
		}
	}

	return data, err
}

func (w *Writer) mergeData(data interface{}) (interface{}, error) {
	var err error

	if merge := w.Cli.String("merge"); merge != "" {
		if padding, dataError := w.readData(merge); dataError == nil {

			for key, val := range padding {
				data.(map[string]interface{})[key] = val
			}
		} else {
			err = dataError
		}
	}

	return data, err
}

func (w *Writer) injectUuid(data interface{}) (interface{}, error) {
	var err error

	if w.Cli.Bool("uuid") {
		if doc, jsonError := json.Marshal(data); jsonError == nil {
			data.(map[string]interface{})["id"] = w.generateUuid(doc)
		} else {
			err = jsonError
		}
	}

	return data, err
}

func (w *Writer) generateUuid(input []byte) string {
	id := uuid.NewSHA1(uuid.NameSpace_DNS, input)
	return id.String()
}

// createOutputDir checks the output flag and creates the
// directory (filemode 0666) on the filesystem if not present
func (w *Writer) createOutputDir() error {
	var err error

	if output := w.Cli.String("output"); output != "" {
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
func (w *Writer) publishData(data interface{}) error {
	var err error
	var id string

	if doc, jsonErr := json.MarshalIndent(data, "", "  "); jsonErr == nil {
		if data.(map[string]interface{})["id"] != nil {
			id = data.(map[string]interface{})["id"].(string)
		} else {
			id = w.generateUuid(doc)
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

	if output := w.Cli.String("output"); output != "" {
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

	if addr := w.Cli.String("address"); addr != "" {
		if uri, uriErr := url.Parse(addr); uriErr == nil {
			client := &http.Client{}

			byteReader := bytes.NewReader(doc)

			if req, httpErr := http.NewRequest(w.Cli.String("http-verb"), uri.String(), byteReader); httpErr == nil {
				req.Header.Set("Content-Type", "application/json")
				resp, reqErr := client.Do(req)
				err = reqErr

				log.Println("HTTP", w.Cli.String("http-verb"), "Response:", resp.Status)
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
		err = w.parseJsonFile(input, &data)
	}

	return data, err
}

func (w *Writer) parseJsonFile(file string, data interface{}) error {
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
