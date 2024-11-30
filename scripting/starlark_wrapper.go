package scripting

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.starlark.net/starlark"
)

func getGlobals(c echo.Context) starlark.StringDict {
	return starlark.StringDict{
		"is_tls": starlark.NewBuiltin("is_tls", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			res := c.IsTLS()
			return starlark.Bool(res), nil
		}),
		"is_websocket": starlark.NewBuiltin("is_websocket", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			res := c.IsWebSocket()
			return starlark.Bool(res), nil
		}),
		"scheme": starlark.NewBuiltin("scheme", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			res := c.Scheme()
			return starlark.String(res), nil
		}),
		"real_ip": starlark.NewBuiltin("real_ip", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			res := c.RealIP()
			return starlark.String(res), nil
		}),
		"path": starlark.NewBuiltin("path", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			res := c.Path()
			return starlark.String(res), nil
		}),
		"set_path": starlark.NewBuiltin("set_path", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var s string
			if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "path", &s); err != nil {
				return nil, err
			}
			c.SetPath(s)
			return nil, nil
		}),
		"param": starlark.NewBuiltin("param", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var s string
			if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "param", &s); err != nil {
				return nil, err
			}
			return starlark.String(c.Param(s)), nil
		}),
		"param_names": starlark.NewBuiltin("param_names", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			res := c.ParamNames()
			vals := starlark.NewList([]starlark.Value{})
			for _, val := range res {
				vals.Append(starlark.String(val))
			}
			return vals, nil
		}),
		"set_param_names": starlark.NewBuiltin("set_param_names", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			names := make([]string, 0)
			for i := 0; i < len(names); i++ {
				name := args[i]
				names = append(names, name.String())
			}
			c.SetParamNames(names...)
			return nil, nil
		}),
		"param_values": starlark.NewBuiltin("param_values", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			vals := c.ParamValues()
			ls := starlark.NewList([]starlark.Value{})
			for _, val := range vals {
				ls.Append(starlark.Value(starlark.String(val)))
			}
			return ls, nil
		}),
		"set_param_values": starlark.NewBuiltin("set_param_values", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			vals := make([]string, 0)
			for i := 0; i < len(vals); i++ {
				val := args[i]
				vals = append(vals, val.String())
			}
			c.SetParamValues(vals...)
			return nil, nil
		}),
		"query_param": starlark.NewBuiltin("query_param", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var s string
			if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "param", &s); err != nil {
				return nil, err
			}
			return starlark.String(c.QueryParam(s)), nil
		}),
		"query_params": starlark.NewBuiltin("query_params", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			params := c.QueryParams()
			ls := starlark.NewList([]starlark.Value{})
			for _, param := range params {
				val := starlark.NewList([]starlark.Value{})
				for _, i := range param {
					val.Append(starlark.Value(starlark.String(i)))
				}
				ls.Append(val)
			}
			return ls, nil
		}),
		"query_string": starlark.NewBuiltin("query_string", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			res := c.QueryString()
			return starlark.String(res), nil
		}),
		"form_value": starlark.NewBuiltin("form_value", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var s string
			if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "value", &s); err != nil {
				return nil, err
			}
			return starlark.String(c.FormValue(s)), nil
		}),
		"form_params": starlark.NewBuiltin("form_params", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			params, err := c.FormParams()
			if err != nil {
				return nil, fmt.Errorf("getting form params: %w", err)
			}
			ls := starlark.NewList([]starlark.Value{})
			for _, param := range params {
				val := starlark.NewList([]starlark.Value{})
				for _, i := range param {
					val.Append(starlark.Value(starlark.String(i)))
				}
				ls.Append(val)
			}
			return ls, nil
		}),
		"get": starlark.NewBuiltin("get", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var s string
			if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "key", &s); err != nil {
				return nil, err
			}
			val, err := interfaceToStarlarkVals(c.Get(s))
			if err != nil {
				return nil, fmt.Errorf("getting starlark values from interface: %w", err)
			}
			return val, nil
		}),
		"set": starlark.NewBuiltin("set", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string
			var val starlark.Value
			if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "key", &key, "value", val); err != nil {
				return nil, err
			}
			v, err := starlarkValToInterface(val)
			if err != nil {
				return nil, err
			}

			c.Set(key, v)
			return nil, nil
		}),
	}
}

func starlarkValToInterface(val starlark.Value) (interface{}, error) {
	switch val.Type() {
	case "bool":
		if v := val.String(); v == "True" {
			return true, nil
		} else {
			return false, nil
		}

	case "int":
		var i int
		err := starlark.AsInt(val, &i)
		if err != nil {
			return nil, fmt.Errorf("getting int: %v", val)
		}
		return i, nil

	case "float":
		f, ok := starlark.AsFloat(val)
		if !ok {
			return nil, fmt.Errorf("getting float64: %v", val)
		}
		return f, nil

	case "string":
		s, ok := starlark.AsString(val)
		if !ok {
			return nil, fmt.Errorf("getting string: %v", val)
		}
		return s, nil

	case "list":
		list, ok := val.(*starlark.List)
		if !ok {
			return nil, fmt.Errorf("getting list: %v", val)
		}
		iter := list.Iterate()

		ls := make([]interface{}, 0)
		var val starlark.Value
		for iter.Next(&val) {
			v, err := starlarkValToInterface(val)
			if err != nil {
				return nil, fmt.Errorf("error")
			}
			ls = append(ls, v)
		}
		return ls, nil

	case "dict":
		dict, ok := val.(*starlark.Dict)
		if !ok {
			return nil, fmt.Errorf("getting dict: %v", val)
		}
		res := make(map[interface{}]interface{})
		for _, key := range dict.Keys() {
			k, err := starlarkValToInterface(key)
			if err != nil {
				return nil, fmt.Errorf("getting interface of key(%v): %w", key, err)
			}

			v, ok, err := dict.Get(key)
			if err != nil {
				return nil, fmt.Errorf("getting value(%v): %w", key, err)
			}
			if !ok {
				return nil, fmt.Errorf("key does not exists: %v", key)
			}
			res[k] = v
		}
		return val, nil

	default:
		return nil, fmt.Errorf("invalid type: %T", val.Type())
	}
}

func interfaceToStarlarkVals(val interface{}) (starlark.Value, error) {
	switch v := val.(type) {
	case string:
		return starlark.String(v), nil
	case int:
		return starlark.MakeInt(v), nil
	case float32:
		return starlark.Float(v), nil
	case float64:
		return starlark.Float(v), nil
	case bool:
		return starlark.Bool(v), nil
	case []interface{}:
		vals := starlark.NewList([]starlark.Value{})
		for _, i := range v {
			val, err := interfaceToStarlarkVals(i)
			if err != nil {
				return nil, fmt.Errorf("getting starlark values from interface: %w", err)
			}
			vals.Append(val)
		}
		return vals, nil
	case map[string]interface{}:
		dict := starlark.NewDict(len(v))
		for key, val := range v {
			sval, err := interfaceToStarlarkVals(val)
			if err != nil {
				return nil, fmt.Errorf("getting starlark value from interface: %v", val)
			}
			dict.SetKey(starlark.String(key), sval)
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("invaild type: %T", v)
	}
}
