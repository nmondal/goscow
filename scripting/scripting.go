package scripting

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/labstack/echo/v4"
)

var jsScripts = make(map[string]*goja.Program)

func runJSProgram(p *goja.Program, context echo.Context) (interface{}, error) {
	vm := goja.New()
	vm.Set("ctx", context)
	v, err := vm.RunProgram(p)
	if err != nil {
		return nil, err
	}
	return v.Export(), nil
}

func doJs(reload bool, filePath string, context echo.Context) (interface{}, error) {
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
	return runJSProgram(p, context)
}

func Execute(reload bool, filePath string, context echo.Context) (interface{}, error) {
	extension := strings.ToLower(filepath.Ext(filePath))
	switch extension {
	case ".js":
		return doJs(reload, filePath, context)
	default:
		panic(fmt.Sprintf("Non Registered extension : '%s'", extension))
	}
}
