package scripting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/labstack/echo/v4"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

var jsScripts = make(map[string]*goja.Program)
var jsRegistry = new(require.Registry) // registry

var starlarkScripts = make(map[string]*starlark.Program)

func runJSProgram(p *goja.Program, context echo.Context) (interface{}, error) {
	vm := goja.New()
	jsRegistry.Enable(vm) // enable it
	console.Enable(vm)    // enable console
	vm.Set("$", context)
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
		data, err := os.ReadFile(filePath) // just pass the file name
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

func runStarlarkProgram(p *starlark.Program, c echo.Context) (interface{}, error) {
	thread := &starlark.Thread{Name: "thread"}
	predeclared := getGlobals(c)

	globals, err := p.Init(thread, predeclared)
	if err != nil {
		return nil, fmt.Errorf("getting starlark globals: %w", err)
	}
	globals.Freeze()
	main := globals["main"]
	v, _ := starlark.Call(thread, main, nil, nil)
	return v, nil
}

func doStarlark(reload bool, filePath string, context echo.Context) (interface{}, error) {
	var p *starlark.Program
	if program, ok := starlarkScripts[filePath]; ok {
		p = program
	} else {
		// compile starlark file
		_, program, err := starlark.SourceProgramOptions(&syntax.FileOptions{}, filePath, nil, func(s string) bool {
			return true
		})

		if err != nil {
			return nil, fmt.Errorf("compiling starlark file: %w", err)
		}

		if !reload {
			starlarkScripts[filePath] = program
		}
		p = program
	}
	return runStarlarkProgram(p, context)
}

func Execute(reload bool, filePath string, context echo.Context) (interface{}, error) {
	extension := strings.ToLower(filepath.Ext(filePath))
	switch extension {
	case ".js":
		return doJs(reload, filePath, context)
	case ".star":
		return doStarlark(reload, filePath, context)
	default:
		panic(fmt.Sprintf("Non Registered extension : '%s'", extension))
	}
}
