package ghostdoc

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/xeipuuv/gojsonschema"
)

type Validator struct {
	Cli *cli.Context
}

func NewValidator(c *cli.Context) *Validator {
	return &Validator{
		Cli: c,
	}
}

/*
  This is fail fast validation. Returns the first error or nil = ok.
*/
func (v *Validator) validate(data map[string]interface{}) error {
	var err error
	if schema := v.Cli.GlobalString("schema"); schema != "" {
		var result *gojsonschema.Result
		schemaLoader := gojsonschema.NewReferenceLoader(schema)
		documentLoader := gojsonschema.NewGoLoader(data)

		result, err = gojsonschema.Validate(schemaLoader, documentLoader)

		if !result.Valid() {
			first := result.Errors()[0]
			err = fmt.Errorf("[Validation error] %v: %v (was %v)", first.Field(), first.Description(), first.Value())
		}
	}
	return err

}
