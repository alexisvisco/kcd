package types

import (
	"encoding"
	"encoding/json"
	"reflect"
	"time"
)

// Unmarshalers is the list of possible unmarshaler kcd support.
var Unmarshalers = []reflect.Type{
	reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem(),
	reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem(),
	reflect.TypeOf((*json.Unmarshaler)(nil)).Elem(),
}

var (
	// UnmarshalerText is the type of TextUnmarshaler
	UnmarshalerText = Unmarshalers[0]

	// BinaryUnmarshaler is the type of BinaryUnmarshaler
	BinaryUnmarshaler = Unmarshalers[1]

	// JSONUnmarshaler is the type of json.Unmarshaler
	JSONUnmarshaler = Unmarshalers[2]
)

// IsImplementingUnmarshaler check if the type t implement one of the possible unmarshalers.
func IsImplementingUnmarshaler(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		t = reflect.New(t).Type()
	}
	for _, u := range Unmarshalers {
		if t.Implements(u) {
			return true
		}
	}

	return false
}

var d = time.Duration(1)

// Custom is the list of custom types supported.
var Custom = []reflect.Type{
	reflect.TypeOf(time.Duration(1)),
	reflect.TypeOf(&d),
}

// IsCustomType check the type is a custom type.
func IsCustomType(t reflect.Type) bool {
	for _, c := range Custom {
		if t.AssignableTo(c) {
			return true
		}
	}
	return false
}

// Native is the list of supported native types.
var Native = map[reflect.Kind]bool{
	reflect.String:  true,
	reflect.Bool:    true,
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
	reflect.Float32: true,
	reflect.Float64: true,
}

// IsNative check if the type t is a native type.
func IsNative(t reflect.Type) bool {
	_, ok := Native[t.Kind()]
	return ok
}

// IsUnmarshallable check if the type t is either a native, custom type or implement an unmarshaler.
func IsUnmarshallable(t reflect.Type) bool {
	return IsNative(t) || IsCustomType(t) || IsImplementingUnmarshaler(t)
}
