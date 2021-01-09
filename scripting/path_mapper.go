package scripting

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	// https://stackoverflow.com/questions/26744873/converting-map-to-struct
	"github.com/mitchellh/mapstructure"
)

// Converts a struct to a map while maintaining the json alias as keys
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

func PathMap(mapping map[string]string, src interface{}, dst interface{}) error {
	vm := goja.New()
	srcMap, err := StructToMap(src)
	if err != nil {
		return err
	}
	dstMap, err := StructToMap(dst)
	if err != nil {
		return err
	}
	vm.Set("$l", srcMap)
	vm.Set("$r", dstMap)

	inx := 1
	for key, value := range mapping {
		scriptInx := fmt.Sprintf("_%d", inx)
		inx++
		scriptBody := fmt.Sprintf("%s = %s;", value, key)
		_, err = vm.RunScript(scriptInx, scriptBody)
		if err != nil {
			return err
		}
	}
	err = mapstructure.Decode(dstMap, dst)
	if err != nil {
		return err
	}
	return nil
}

type MyStruct2 struct {
	X int64
	Y int64
}

type MyStruct1 struct {
	Name string
	Age  int64
	MS   MyStruct2
}

func TestIt() {
	src := MyStruct1{
		Name: "Yo!",
		Age:  42,
		MS: MyStruct2{
			X: 12,
			Y: 13,
		},
	}
	fmt.Printf("SRC : %v\n", src)
	dst := MyStruct1{MS: MyStruct2{}}
	fmt.Printf("DST Before : %v\n", dst)
	mapping := map[string]string{
		"$l.Name": "$r.Name",
		"$l.Age":  "$r.Age",
		"$l.MS.X": "$r.MS.X",
		"$l.MS.Y": "$r.MS.Y",
	}
	// notice this, one has to provide the pointer...
	e := PathMap(mapping, src, &dst)
	if e != nil {
		panic(e)
	}
	fmt.Printf("DST After : %v\n", dst)
}
