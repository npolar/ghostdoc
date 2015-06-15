package ghostdoc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/npolar/ghostdoc/context"
)

type publisher struct {
	context   context.GhostContext
	documents []map[string]interface{}
	sem       chan int
}

func newPublisher(c context.GhostContext) *publisher {
	return &publisher{context: c, documents: make([]map[string]interface{}, 0), sem: make(chan int, 1)}
}

func (p *publisher) add(data map[string]interface{}) {
	p.sem <- 1
	p.documents = append(p.documents, data)
	<-p.sem
}

func (p *publisher) send() error {
	var err error
	bulkSize := p.context.GlobalInt("bulk")
	p.sem <- 1
	for {
		if bulkSize > len(p.documents) {
			bulkSize = len(p.documents)
		}
		chunk := p.documents[:bulkSize]
		p.documents = p.documents[bulkSize:]

		if docs, jsonErr := json.MarshalIndent(chunk, "", "  "); jsonErr == nil {
			err = p.httpRequest(docs)
		} else {
			err = errors.New("publishData: " + jsonErr.Error())
		}

		if err != nil {
			return err
		}

		if len(p.documents) == 0 {
			break
		}
	}
	<-p.sem
	return err
}

// httpRequest performs the request defined in http-verb against
// the configured address. Default operation is POST
func (p *publisher) httpRequest(docs []byte) error {
	var err error
	if addr := p.context.GlobalString("address"); addr != "" {
		var verb = p.context.GlobalString("http-verb")
		if uri, uriErr := url.Parse(addr); uriErr == nil {
			client := &http.Client{}

			byteReader := bytes.NewReader(docs)

			if req, httpErr := http.NewRequest(verb, uri.String(), byteReader); httpErr == nil {
				req.Header.Set("Content-Type", "application/json")
				var resp *http.Response

				if resp, err = client.Do(req); err == nil {
					defer resp.Body.Close()
					log.Debug("HTTP", verb, "Response:", resp.Status)
				} else {
					log.Error("HTTP Error", verb, err.Error())
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
