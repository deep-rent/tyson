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

// Package tyson provides a wrapper for navigating hierarchies of unknown or
// dynamic JSON objects parsed into map[string]any. It aims to hide the messy
// type and existence checks that usually occur when dealing with such
// structures.
package tyson

// A Mapper converts type S into type T. If the conversion succeeded, ok will be
// true; if not, ok will be false. This is essentially a generic abstraction of
// the "comma ok" idiom used in Go type casts.
type Mapper[S any, T any] func(v S) (w T, ok bool)

/* Default mappers */

func AsArray(v any) (w []any, ok bool)   { w, ok = v.([]any); return }
func AsBool(v any) (w bool, ok bool)     { w, ok = v.(bool); return }
func AsFloat(v any) (w float64, ok bool) { w, ok = v.(float64); return }
func AsInt(v float64) (w int64, ok bool) { w = int64(v); return w, v == float64(w) }
func AsObject(v any) (w Object, ok bool) { w, ok = v.(map[string]any); return }
func AsString(v any) (w string, ok bool) { w, ok = v.(string); return }

// All "lifts" m to convert a slice of item type S into a slice of item type T.
// The resulting [Mapper] indicates ok if and only if m was successfully
// applied to each element of the input slice.
func All[S any, T any](m Mapper[S, T]) Mapper[[]S, []T] {
	return func(v []S) (w []T, ok bool) {
		w = make([]T, len(v))
		for i, x := range v {
			if y, ok := m(x); ok {
				w[i] = y
			} else {
				return nil, false
			}
		}
		return w, true
	}
}

// One returns a [Mapper] that applies the given mappers one after another
// until the first to succeed. It fails if none of the mappers succeed.
func One[S any, T any](mappers ...Mapper[S, T]) Mapper[S, T] {
	return func(v S) (w T, ok bool) {
		for _, m := range mappers {
			w, ok = m(v)
			if ok {
				break
			}
		}
		return
	}
}

// A Node is a container which may or may not contain a value of type T. If a
// value is present, Empty returns false and Value will return the value.
// Additional methods that depend on the presence or absence of a contained
// value are provided, such as Or which returns a default value if this Node is
// empty.
//
// Use [EmptyNode] or [ValueNode] to create empty or nonempty Nodes
// respectively.
type Node[T any] interface {
	// Value returns the contained value if present, or else the
	// zero value of T.
	Value() T
	// Empty returns true if no value is present, otherwise false.
	Empty() bool
	// Or returns the contained value if present, otherwise v.
	Or(v T) T
	// OrGet returns the contained value if present, or else the
	// result produced by f.
	OrGet(f func() T) T
}

// EmptyNode returns an empty [Node] of type T.
func EmptyNode[T any]() Node[T] {
	return empty[T]{}
}

// ValueNode returns a nonempty [Node] containing v.
func ValueNode[T any](v T) Node[T] {
	return value[T]{value: v}
}

type empty[T any] struct{}

func (empty[T]) Value() T           { return *new(T) }
func (empty[T]) Empty() bool        { return true }
func (empty[T]) Or(v T) T           { return v }
func (empty[T]) OrGet(f func() T) T { return f() }

type value[T any] struct{ value T }

func (v value[T]) Value() T         { return v.value }
func (v value[T]) Empty() bool      { return false }
func (v value[T]) Or(T) T           { return v.value }
func (v value[T]) OrGet(func() T) T { return v.value }

// Map applies m to convert the value contained in n. If n is empty, the
// resulting [Node] will be empty as well.
func Map[S any, T any](n Node[S], m Mapper[S, T]) Node[T] {
	if !n.Empty() {
		if v, ok := m(n.Value()); ok {
			return ValueNode(v)
		}
	}
	return EmptyNode[T]()
}

// An Object represents the result of parsing an unknown or dynamic JSON object.
// Provided getters can be used to navigate the object hierarchy in a type-safe
// manner.
type Object map[string]any

// Get follows the given hierarchy of keys to locate a target value within the
// underlying JSON structure. The returned [Node] is empty if some key does
// not exist, or else contains the target value. If no key is passed, the
// returned [Node] contains this [Object].
//
// The following example fetches the nested value of "c" from the parsed JSON
// object:
//
//	var o tyson.Object
//	_ = json.Unmarshal([]byte(`{"a":{"b":{"c":"d"}}}`), &o)
//	fmt.Print(o.Get("a", "b", "c").Value()) // prints "d"
func (o Object) Get(keys ...string) Node[any] {
	switch len(keys) {
	case 0:
		return ValueNode[any](o)
	case 1:
		return o.get(keys[0])
	default:
		n := ValueNode(o)
		i := 0
		for i < len(keys)-1 {
			n = Map(n.Value().get(keys[i]), AsObject)
			if n.Empty() {
				break
			}
			i++
		}
		return n.Value().get(keys[i])
	}
}

func (o Object) get(k string) Node[any] {
	if v, ok := o[k]; ok {
		return ValueNode(v)
	} else {
		return EmptyNode[any]()
	}
}

// Has returns true if key k is present in this Object, otherwise false.
func (o Object) Has(k string) bool { _, ok := o[k]; return ok }

// Set assigns value v to key k.
func (o Object) Set(k string, v any) { o[k] = v }

// Remove deletes the value with key k from this Object.
func (o Object) Remove(k string) { delete(o, k) }

/* Basic type getters */

// GetArray follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON array. Otherwise, the
// returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetArray(keys ...string) Node[[]any] {
	return Map(o.Get(keys...), AsArray)
}

// GetBool follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON boolean. Otherwise, the
// returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetBool(keys ...string) Node[bool] {
	return Map(o.Get(keys...), AsBool)
}

// GetFloat follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON number. Otherwise, the
// returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetFloat(keys ...string) Node[float64] {
	return Map(o.Get(keys...), AsFloat)
}

// GetInt follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not an integral JSON number.
// Otherwise, the returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetInt(keys ...string) Node[int64] {
	return Map(o.GetFloat(keys...), AsInt)
}

// GetObject follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON object. Otherwise, the
// returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetObject(keys ...string) Node[Object] {
	return Map(o.Get(keys...), AsObject)
}

// GetString follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON string. Otherwise, the
// returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetString(keys ...string) Node[string] {
	return Map(o.Get(keys...), AsString)
}

/* Array type getters */

// GetArrays follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a two-dimensional JSON array.
// Otherwise, the returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetArrays(keys ...string) Node[[][]any] {
	return Map(o.GetArray(keys...), All(AsArray))
}

// GetBools follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON array of only booleans.
// Otherwise, the returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetBools(keys ...string) Node[[]bool] {
	return Map(o.GetArray(keys...), All(AsBool))
}

// GetFloats follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON array of only numbers.
// Otherwise, the returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetFloats(keys ...string) Node[[]float64] {
	return Map(o.GetArray(keys...), All(AsFloat))
}

// GetInts follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON array of only integral
// numbers. Otherwise, the returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetInts(keys ...string) Node[[]int64] {
	return Map(o.GetFloats(keys...), All(AsInt))
}

// GetObjects follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON array of only objects.
// Otherwise, the returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetObjects(keys ...string) Node[[]Object] {
	return Map(o.GetArray(keys...), All(AsObject))
}

// GetStrings follows the given hierarchy of keys to fetch a target value from
// the underlying JSON structure. The returned [Node] is empty if some key
// does not exist, or if the target value is not a JSON array of only strings.
// Otherwise, the returned [Node] contains the target value.
//
// See [Object.Get] for an example on how to fetch nested JSON values.
func (o Object) GetStrings(keys ...string) Node[[]string] {
	return Map(o.GetArray(keys...), All(AsString))
}
