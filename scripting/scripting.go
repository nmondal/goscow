package scripting

import (
	"github.com/dop251/goja"
	"io/ioutil"
)

var jsScripts = make(map[string]*goja.Program)

func runProgram(p *goja.Program, args map[string]interface{}) (interface{}, error) {
	vm := goja.New()
	for varName, value := range args {
		vm.Set(varName, value)
	}
	v, err := vm.RunProgram(p)
	if err != nil {
		return nil, err
	}
	return v.Export(), nil
}

func JS(filePath string, args map[string]interface{}) (interface{}, error) {
	var p *goja.Program
	if program, ok := jsScripts[filePath]; ok {
		p = program
	} else {
		data, err := ioutil.ReadFile(filePath) // just pass the file name
		if err != nil {
			return nil, err
		}
		scriptString := string(data)
		program, err := goja.Compile(filePath, scriptString, false)
		if err != nil {
			return nil, err
		}
		p = program
	}
	return runProgram(p, args)
}
