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
	"reflect"
	"strconv"
	"testing"

	"github.com/deep-rent/tyson"
)

func TestAll(t *testing.T) {
	m := tyson.All(func(v int) (string, bool) {
		return strconv.Itoa(v), true
	})

	exp := []string{"1", "2", "3"}
	act, ok := m([]int{1, 2, 3})

	if !ok {
		t.Fatalf("ok was %t", ok)
	}

	if !reflect.DeepEqual(act, exp) {
		t.Errorf("want %v, was %v", exp, act)
	}
}

func TestAll_Fail(t *testing.T) {
	m := tyson.All(func(v int) (string, bool) {
		return "", false
	})

	_, ok := m([]int{1, 2, 3})
	if ok {
		t.Fatalf("ok was %t", ok)
	}
}

func TestOne(t *testing.T) {
	f := func(v int) (string, bool) {
		return "", false
	}
	s := func(v int) (string, bool) {
		return strconv.Itoa(v), true
	}

	m := tyson.One(f, s, f)

	exp := "1"
	act, ok := m(1)

	if !ok {
		t.Fatalf("ok was %t", ok)
	}

	if act != exp {
		t.Errorf("want %v, was %v", exp, act)
	}
}

func TestOne_Fail(t *testing.T) {
	f := func(v int) (string, bool) {
		return "", false
	}

	m := tyson.One(f, f, f)

	_, ok := m(1)
	if ok {
		t.Fatalf("ok was %t", ok)
	}
}

func TestMap(t *testing.T) {
	n := tyson.ValueNode(123)
	v := tyson.Map(n, func(v int) (string, bool) {
		return strconv.Itoa(v), true
	})

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := "123"
	act := v.Value()

	if act != exp {
		t.Errorf("was %q, want %q", act, exp)
	}
}

func TestMap_Fail(t *testing.T) {
	n := tyson.ValueNode(123)
	v := tyson.Map(n, func(v int) (string, bool) {
		return "", false
	})

	if !v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}
}

func TestObject_Get(t *testing.T) {
	o := make(tyson.Object)

	exp := o
	act := o.Get().Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_Get_Single(t *testing.T) {
	o := make(tyson.Object)
	o.Set("foo", "bar")

	exp := "bar"
	act := o.Get("foo").Value()

	if exp != act {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_Get_Nested(t *testing.T) {
	o := make(tyson.Object)
	o.Set("foo", map[string]any{"bar": "baz"})

	exp := "baz"
	act := o.Get("foo", "bar").Value()

	if exp != act {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_Get_Nested_Fail(t *testing.T) {
	o := make(tyson.Object)
	o.Set("foo", map[string]any{"bar": "baz"})

	v := o.Get("foo", "bar", "baz")

	if !v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}
}

func TestObject_GetArray(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", []any{1, 2})

	v := o.GetArray("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := []any{1, 2}
	act := v.Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_GetBool(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", true)

	v := o.GetBool("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := true
	act := v.Value()

	if act != exp {
		t.Fatalf("was %t, want %t", act, exp)
	}
}

func TestObject_GetFloat(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", float64(12.34))

	v := o.GetFloat("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := float64(12.34)
	act := v.Value()

	if act != exp {
		t.Fatalf("was %f, want %f", act, exp)
	}
}

func TestObject_GetInt(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", float64(10))

	v := o.GetInt("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := int64(10)
	act := v.Value()

	if act != exp {
		t.Fatalf("was %d, want %d", act, exp)
	}
}

func TestObject_GetObject(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", map[string]any{"foo": "bar"})

	v := o.GetObject("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := tyson.Object{"foo": "bar"}
	act := v.Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_GetString(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", "abc")

	v := o.GetString("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := "abc"
	act := v.Value()

	if act != exp {
		t.Fatalf("was %q, want %q", act, exp)
	}
}

func TestObject_GetArrays(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", []any{[]any{}, []any{}})

	v := o.GetArrays("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := [][]any{{}, {}}
	act := v.Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_GetBools(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", []any{true, false})

	v := o.GetBools("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := []bool{true, false}
	act := v.Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_GetFloats(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", []any{float64(12.34), float64(56.78)})

	v := o.GetFloats("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := []float64{12.34, 56.78}
	act := v.Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_GetInts(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", []any{float64(10.00), float64(20.00)})

	v := o.GetInts("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := []int64{10, 20}
	act := v.Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_GetObjects(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", []any{map[string]any{}, map[string]any{}})

	v := o.GetObjects("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := []tyson.Object{{}, {}}
	act := v.Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}

func TestObject_GetStrings(t *testing.T) {
	o := make(tyson.Object)
	o.Set("v", []any{"abc", "def"})

	v := o.GetStrings("v")

	if v.Empty() {
		t.Fatalf("empty was %t", v.Empty())
	}

	exp := []string{"abc", "def"}
	act := v.Value()

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("was %#v, want %#v", act, exp)
	}
}
