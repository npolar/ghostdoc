package ghostdoc

import (
	"fmt"

	"github.com/npolar/ghostdoc/context"
	"github.com/xeipuuv/gojsonschema"
)

// Validator typedef
type Validator struct {
	context context.GhostContext
	schema  gojsonschema.JSONLoader
}

// NewValidator factory
func NewValidator(c context.GhostContext) *Validator {
	var jsonRef gojsonschema.JSONLoader
	if s := c.GlobalString("schema"); s != "" {
		jsonRef = gojsonschema.NewReferenceLoader(s)
	}
	return &Validator{
		context: c,
		schema:  jsonRef,
	}
}

/*
  This is fail fast validation. Returns the first error or nil = ok.
*/
func (v *Validator) validate(data map[string]interface{}) error {
	var err error
	if v.schema != nil {
		var result *gojsonschema.Result
		schemaLoader := v.schema
		documentLoader := gojsonschema.NewGoLoader(data)

		if result, err = gojsonschema.Validate(schemaLoader, documentLoader); err == nil {
			if !result.Valid() {
				first := result.Errors()[0]
				err = fmt.Errorf("[Validation error] %v: %v (was %v)", first.Field(), first.Description(), first.Value())
			}
		}
	}
	return err

}
