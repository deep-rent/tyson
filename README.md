# deep-rent/tyson

![Logo](https://raw.githubusercontent.com/deep-rent/tyson/master/logo.png)

Tyson is a tiny Go library that helps with navigating hierarchies of unknown or dynamic JSON structures parsed into `map[string]any`. It aims to hide the messy type and existence checks that usually occur when dealing with such structures.

## Installation

Download the libary using `go get`:

```
go get github.com/deep-rent/tyson@latest
```

Add the following import to your project:

```go
import (
    "github.com/deep-rent/tyson"
)
```

## Usage

Start by parsing your JSON into `tyson.Object`. This type is simply an alias for `map[string]any`:

```go
var o tyson.Object
_ = json.Unmarshal([]byte(`
{
    "obj": {
        "num": 42,
        "str": "abc"
    },
    "arr": [true, false]
}
`), &o)
```

Using the `GetXXX()` methods, you can fetch specific values from the wrapped structure in a type-safe manner. Getters are available for Go equivalents of all JSON types including arrays. Each getter can take multiple keys to target nested values:

```go
fmt.Printf("num: %d\n", o.GetInt("obj", "num").Value())
fmt.Printf("str: %q\n", o.GetString("obj", "str").Value())
fmt.Printf("arr: %v\n", o.GetBools("arr").Value())
```

Output:

```
num: 42
str: "abc"
arr: [true false]
```

Note that the getters don't directly return the target values but rather instances of the `tyson.Node` interface. This is important because the target value might be of an unexpected type, or might not exist at all. In that case, the returned node is `Empty()`, and a default value can be provided using `Or()` and `OrGet()`:

```go
xyz := o.Get("xyz")
str := o.GetString("obj", "num")

fmt.Printf(`key "xyz" does not exist: %t`+"\n", xyz.Empty())
fmt.Printf(`val "num" isn't a string: %t`+"\n", str.Empty())

fmt.Printf("xyz: %v\n", xyz.Or("def"))
fmt.Printf("str: %s\n", str.OrGet(func() string { return "def" }))
```

Output:

```
key "xyz" does not exist: true
val "num" isn't a string: true
xyz: def
str: def
```



## Remarks

- Parsing JSON into `map[string]any` is known to be slow. If you deal with large JSON objects, or target only a small subset of the entire structure, you should look for libraries that implement specialized parsers.
- The `GetInt` and `GetInts` methods of `Object` require the target value to actually be an integral number. For example, `42` and `42.00` are valid integers, but `42.01` causes the getters to return an empty `Node.`

## License

Licensed under the Apache 2.0 License. For the full copyright and licensing information, please view the `LICENSE` file that was distributed with this source code.