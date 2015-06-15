package ghostdoc

import (
	"errors"
	"io/ioutil"

	"github.com/npolar/ghostdoc/context"
	"github.com/robertkrimen/otto"
)

// Js typedef
type Js struct {
	context   context.GhostContext
	vm        *otto.Otto
	functions *otto.Object
}

// NewJs factory
func NewJs(c context.GhostContext) *Js {
	js := &Js{
		context: c,
		vm:      otto.New(),
	}

	if flag := js.context.GlobalString("js"); flag != "" {
		if err := js.runCode(flag); err != nil {
			panic("JavaScript does not compile")
		} else {
			if fns, err := js.vm.Get("functions"); err == nil {
				fnObject := fns.Object()
				if fnObject != nil {
					js.functions = fnObject
				} else {
					panic("Could not find \"functions\" Object in script")
				}
			}
		}
	}

	return js
}

func (js *Js) copy() *Js {
	vm := js.vm.Copy()
	fns, _ := vm.Get("functions")
	return &Js{
		context:   js.context,
		vm:        js.vm.Copy(),
		functions: fns.Object(),
	}
}

func (js *Js) runJs(data map[string]interface{}) (map[string]interface{}, error) {
	var err error
	if flag := js.context.GlobalString("js"); flag != "" {
		data, err = js.runFunctions(data)
	}
	if err != nil {
		err = errors.New("[JavaScript Error] " + err.Error())
	}
	return data, err
}

func (js *Js) runFunctions(data map[string]interface{}) (map[string]interface{}, error) {
	var err error
	var fns otto.Value
	if fns, err = js.vm.Get("functions"); err == nil {
		fnObject := fns.Object()
		if fnObject != nil {
			for _, fnName := range fnObject.Keys() {
				data, err = js.runFunction(fnObject, fnName, data)
			}
		} else {
			err = errors.New("Could not find \"functions\" Object in script")
		}

	}
	return data, err
}

func (js *Js) runFunction(fnObject *otto.Object, fnName string, data map[string]interface{}) (map[string]interface{}, error) {
	var err error
	var response otto.Value
	if response, err = fnObject.Call(fnName, js.makeOttoObject(data)); err == nil {
		export, _ := response.Export()
		if export != nil {
			data = export.(map[string]interface{})
		} else {
			data = nil
		}
	}
	return data, err
}

func (js *Js) runCode(flag string) error {
	var err error
	var script *otto.Script
	code := js.readCode(flag)
	if script, err = js.vm.Compile("", code); err == nil {
		js.vm.Run(script)
	}
	return err
}

func (js *Js) readCode(flag string) string {
	code := flag
	if raw, err := ioutil.ReadFile(flag); err == nil {
		code = string(raw)
	}
	return code
}

func (js *Js) makeOttoObject(object map[string]interface{}) *otto.Object {
	ottoObject, _ := js.vm.Object(`({})`)
	for key, value := range object {
		ottoObject.Set(key, value)
	}
	return ottoObject
}
