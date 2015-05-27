package ghostdoc

import (
	"errors"
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/robertkrimen/otto"
)

type Js struct {
	Cli *cli.Context
	vm  *otto.Otto
}

func NewJs(c *cli.Context) *Js {
	return &Js{
		Cli: c,
		vm:  otto.New(),
	}
}

func (js *Js) runJs(data map[string]interface{}) (map[string]interface{}, error) {
	var err error
	if flag := js.Cli.GlobalString("js"); flag != "" {
		if err = js.runCode(flag); err == nil {
			data, err = js.runFunctions(data)
		}
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
		data = export.(map[string]interface{})
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
