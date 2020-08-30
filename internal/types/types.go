package types

import (
	"encoding"
	"encoding/json"
	"reflect"
	"time"
)

var Unmarshallers = []reflect.Type{
	reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem(),
	reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem(),
	reflect.TypeOf((*json.Unmarshaler)(nil)).Elem(),
}

var (
	UnmarshallerText   = Unmarshallers[0]
	BinaryUnmarshaller = Unmarshallers[1]
	JsonUnmarshaller   = Unmarshallers[2]
)

func IsImplementingUnmarshaller(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		t = reflect.New(t).Type()
	}
	for _, u := range Unmarshallers {
		if t.Implements(u) {
			return true
		}
	}

	return false
}

var d = time.Duration(1)

var Custom = []reflect.Type{
	reflect.TypeOf(time.Duration(1)),
	reflect.TypeOf(&d),
}

func IsCustomType(t reflect.Type) bool {
	for _, c := range Custom {
		if t.AssignableTo(c) {
			return true
		}
	}
	return false
}

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

func IsNative(t reflect.Type) bool {
	_, ok := Native[t.Kind()]
	return ok
}

func IsUnmarshallable(t reflect.Type) bool {
	return IsNative(t) || IsCustomType(t) || IsImplementingUnmarshaller(t)
}
