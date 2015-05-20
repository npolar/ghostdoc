package ghostdoc

//import (
//	"encoding/json"
//	"io/ioutil"
//	"log"
//	"regexp"
//	"strconv"
//)
//
//type Parser struct{}
//
//// parseFilename handles meta data extraction from filenames. It reads the pattern
//// file specified with the --name-pattern argument and parses the filename according
//func (p *Parser) parseFileName(fname string, doc interface{}) (interface{}, error) {
//	var err error
//
//	if pat := p.Cli.String("name-pattern"); pat != "" {
//		var pattern = make(map[string]interface{})
//
//		pdoc, err := ioutil.ReadFile(pat)
//		err = json.Unmarshal(pdoc, &pattern)
//
//		if err == nil {
//			if pRgx := regexp.MustCompile(pattern["pattern"].(string)); pRgx.MatchString(fname) {
//				matches := pRgx.FindStringSubmatch(fname)
//				log.Println(matches)
//				output := pattern["output"].(string)
//
//				// Replace the %<count> indicators in the output with the matching capture
//				for i, match := range matches {
//					rxp := regexp.MustCompile("(%" + strconv.Itoa(i) + ")")
//					output = rxp.ReplaceAllString(output, match)
//				}
//
//				doc.(map[string]interface{})[pattern["key"].(string)] = output
//			}
//		}
//	}
//
//	return doc, err
//}
