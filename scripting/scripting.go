package scripting

import (
	"fmt"
	"github.com/dop251/goja"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var jsScripts = make(map[string]*goja.Program)

func runJSProgram(p *goja.Program, args map[string]interface{}) (interface{}, error) {
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

func doJs(reload bool, filePath string, args map[string]interface{}) (interface{}, error) {
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
		if !reload {
			jsScripts[filePath] = program
		}
		p = program
	}
	return runJSProgram(p, args)
}

func Execute(reload bool, filePath string, args map[string]interface{}) (interface{}, error) {
	extension := strings.ToLower(filepath.Ext(filePath))
	switch extension {
	case ".js":
		return doJs(reload, filePath, args)
	default:
		panic(fmt.Sprintf("Non Registered extension : '%s'", extension))
	}
}
