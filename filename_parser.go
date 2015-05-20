package ghostdoc

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/codegangsta/cli"
)

// parseFilename handles meta data extraction from filenames. It reads the pattern
// file specified with the --name-pattern argument and parses the filename according
func parseFileName(Cli *cli.Context, fname string, doc interface{}) (interface{}, error) {
	var err error
	if pat := Cli.GlobalString("name-pattern"); pat != "" {
		var pattern = make(map[string]interface{})
		var pdoc []byte
		if pdoc, err = ioutil.ReadFile(pat); err != nil {
			pdoc = []byte(pat)
		}

		if err = json.Unmarshal(pdoc, &pattern); err == nil {
			if pRgx := regexp.MustCompile(pattern["pattern"].(string)); pRgx.MatchString(fname) {
				matches := pRgx.FindStringSubmatch(fname)
				output := pattern["output"].(string)

				// Replace the %<count> indicators in the output with the matching capture
				for i, match := range matches {
					rxp := regexp.MustCompile("(%" + strconv.Itoa(i) + ")")
					output = rxp.ReplaceAllString(output, match)
				}

				doc.(map[string]interface{})[pattern["key"].(string)] = output
			}
		}
	}

	return doc, err
}
