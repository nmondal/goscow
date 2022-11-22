package scripting

import (
	"github.com/itchyny/gojq"
)

/*
https://stedolan.github.io/jq/manual/
JQ is used to format JSON
*/

func Jq(input interface{}, jqScript string) ([]interface{}, error) {
	query, err := gojq.Parse(jqScript)
	if err != nil {
		return nil, err
	}
	iter := query.Run(input) // or query.RunWithContext
	var result []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func JqSingle(input interface{}, jqScript string) (interface{}, error) {
	result, err := Jq(input, jqScript)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return result, nil
	}
	return result[0], nil
}
