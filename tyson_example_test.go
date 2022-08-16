/*
Copyright (c) 2022 deep.rent GmbH (https://deep.rent)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tyson_test

import (
	"encoding/json"
	"fmt"

	"github.com/deep-rent/tyson"
)

// This example showcases the retrieval of various types from a parsed JSON
// object.
func Example_type_getters() {
	var o tyson.Object
	_ = json.Unmarshal([]byte(`
	{
		"obj": {
			"num": 12.34,
			"int": 12345,
			"str": "abc"
		},
		"arr": [true, false]
	}
	`), &o)

	fmt.Printf("num: %.2f\n", o.GetFloat("obj", "num").Value())
	fmt.Printf("int: %d\n", o.GetInt("obj", "int").Value())
	fmt.Printf("str: %q\n", o.GetString("obj", "str").Value())
	fmt.Printf("arr: %v\n", o.GetBools("arr").Value())

	// Output:
	// num: 12.34
	// int: 12345
	// str: "abc"
	// arr: [true false]
}

// This example explains how to return a default value for target values that
// either don't exist or have a wrong type.
func Example_default_values() {
	var o tyson.Object
	_ = json.Unmarshal([]byte(`{ "num": 12.34 }`), &o)

	xyz := o.Get("xyz")
	str := o.GetString("num")

	fmt.Printf(`key "xyz" does not exist: %t`+"\n", xyz.Empty())
	fmt.Printf(`val "num" isn't a string: %t`+"\n", str.Empty())

	fmt.Printf("xyz: %v\n", xyz.Or("def"))
	fmt.Printf("str: %s\n", str.OrGet(func() string { return "def" }))

	// Output:
	// key "xyz" does not exist: true
	// val "num" isn't a string: true
	// xyz: def
	// str: def
}
