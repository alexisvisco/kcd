package decoder

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/expectedsh/errors"

	"github.com/expectedsh/kcd/internal/cache"
	"github.com/expectedsh/kcd/internal/kcderr"
	"github.com/expectedsh/kcd/internal/types"
)

type setterContext struct {
	metadata  cache.FieldMetadata
	path, tag string
	value     interface{} // value is either []string or interface{}
}

type fieldSetter struct {
	field     reflect.Value
	value     interface{}
	metadata  cache.FieldMetadata
	errFields map[string]interface{}
}

func newFieldSetter(field reflect.Value, setterCtx setterContext) *fieldSetter {
	fs := &fieldSetter{field: field, value: setterCtx.value, metadata: setterCtx.metadata}

	fs.errFields = map[string]interface{}{
		"value":      fmt.Sprint(setterCtx.value),
		"value-type": reflect.TypeOf(setterCtx.value).String(),
		"field-type": field.Type().String(),
		"path":       setterCtx.path,
		"tag":        setterCtx.tag,
	}

	return fs
}
func (f fieldSetter) set() error {
	if f.field.Type().AssignableTo(reflect.TypeOf(f.value)) {
		f.field.Set(reflect.ValueOf(f.value))
		return nil
	}

	list, ok := f.value.([]string)
	if !ok {
		switch t := f.value.(type) {
		case string:
			list = []string{t}
		case []byte:
			list = []string{string(t)}
		default:
			return errors.NewWithKind(kcderr.InputCritical, "incompatible type").WithFields(f.errFields)
		}
	}

	if len(list) == 0 {
		return nil
	}

	isPtr := f.field.Kind() == reflect.Ptr

	if f.metadata.ArrayOrSlice {
		return f.setForArrayOrSlice(isPtr, list)
	}

	return f.setForNormalType(list[0], isPtr)
}

func (f fieldSetter) setForArrayOrSlice(ptr bool, list []string) error {
	var (
		element reflect.Value
		array   = false
	)

	isTypePtr := false
	elemType := f.field.Type()
	array = f.field.Type().Kind() == reflect.Array

	if ptr {
		array = f.field.Type().Elem().Kind() == reflect.Array
		elemType = elemType.Elem()
		isTypePtr = f.field.Type().Elem().Elem().Kind() == reflect.Ptr
	} else {
		isTypePtr = f.field.Type().Elem().Kind() == reflect.Ptr
	}

	if array {
		element = reflect.New(elemType)
	} else {
		element = reflect.MakeSlice(elemType, 0, len(list))
	}

	addToElem := func(index int, i reflect.Value) {
		if array {
			element.Elem().Index(index).Set(i)
		} else {
			element = reflect.Append(element, i)
		}
	}

	for i, val := range list {
		if types.IsCustomType(f.metadata.Type) {
			native, err := f.makeCustomType(val, isTypePtr)
			if err != nil {
				return err.WithField("value-index", i)
			}

			addToElem(i, native)
		} else if types.IsNative(f.metadata.Type) {
			native, err := f.makeNative(val, isTypePtr)
			if err != nil {
				return err.WithField("value-index", i)
			}

			addToElem(i, native)
		} else if types.IsImplementingUnmarshaller(f.metadata.Type) {
			withUnmarshaller, err := f.makeWithUnmarshaller(val)
			if err != nil {
				return err.WithField("value-index", i)
			}

			addToElem(i, withUnmarshaller)
		} else {
			return errors.NewWithKind(kcderr.InputCritical, "type is not native, unmarshaller or custom type").
				WithFields(f.errFields).
				WithField("value-index", i)
		}
	}

	setToField := func(value reflect.Value) {
		if element.Kind() == reflect.Ptr && element.Elem().Kind() == reflect.Array {
			value.Set(element.Elem())
		} else {
			value.Set(element)
		}
	}

	if ptr {
		v := reflect.New(f.field.Type().Elem())
		setToField(v.Elem())
		f.field.Set(v)
	} else {
		setToField(f.field)
	}

	return nil
}

func (f fieldSetter) setForNormalType(str string, ptr bool) error {
	if types.IsCustomType(f.metadata.Type) {
		customType, err := f.makeCustomType(str, ptr)
		if err != nil {
			return err
		}

		f.field.Set(customType)
	} else if types.IsNative(f.metadata.Type) {
		native, err := f.makeNative(str, ptr)
		if err != nil {
			return err
		}

		f.field.Set(native)
	} else if types.IsImplementingUnmarshaller(f.metadata.Type) {
		withUnmarshaller, err := f.makeWithUnmarshaller(str)
		if err != nil {
			return err
		}

		f.field.Set(withUnmarshaller)
	} else {
		return errors.NewWithKind(kcderr.InputCritical, "type is not native, unmarshaller or custom type").
			WithFields(f.errFields)
	}

	return nil
}

func (f fieldSetter) makeNative(str string, ptr bool) (reflect.Value, *errors.Error) {
	el := reflect.New(f.metadata.Type)

	switch f.metadata.Type.Kind() {
	case reflect.String:
		el.Elem().SetString(str)

		if ptr {
			return el, nil
		}
		return el.Elem(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, f.metadata.Type.Bits())
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "invalid integer").
				WithKind(kcderr.Input).
				WithFields(f.errFields)
		}

		el.Elem().SetInt(i)

		if ptr {
			return el, nil
		}
		return el.Elem(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(str, 10, f.metadata.Type.Bits())
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "invalid positive integer").
				WithKind(kcderr.Input).
				WithFields(f.errFields)
		}

		el.Elem().SetUint(i)

		if ptr {
			return el, nil
		}
		return el.Elem(), nil
	case reflect.Bool:
		i, err := strconv.ParseBool(str)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "invalid boolean").
				WithKind(kcderr.Input).
				WithFields(f.errFields)
		}

		el.Elem().SetBool(i)

		if ptr {
			return el, nil
		}
		return el.Elem(), nil
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(str, f.metadata.Type.Bits())
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "invalid floating number").
				WithKind(kcderr.Input).
				WithFields(f.errFields)
		}

		el.Elem().SetFloat(i)

		if ptr {
			return el, nil
		}
		return el.Elem(), nil
	}

	return reflect.Value{}, errors.NewWithKind(kcderr.InputCritical, "an error occur with this getValueFromHttp").
		WithFields(f.errFields)
}

func (f fieldSetter) makeWithUnmarshaller(str string) (reflect.Value, *errors.Error) {
	var el reflect.Value
	if f.metadata.Type.Kind() == reflect.Ptr {
		el = reflect.New(f.metadata.Type.Elem())
	} else {
		el = reflect.New(f.metadata.Type)
	}

	if el.Type().Implements(types.UnmarshallerText) {

		t := el.Interface().(encoding.TextUnmarshaler)

		err := t.UnmarshalText([]byte(str))
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "unable to unmarshal from binary format").
				WithKind(kcderr.Input).
				WithFields(f.errFields)
		}

		return el, nil
	}

	if el.Type().Implements(types.JsonUnmarshaller) {
		t := el.Interface().(json.Unmarshaler)
		err := t.UnmarshalJSON([]byte(str))
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "unable to unmarshal from json format").
				WithKind(kcderr.Input).
				WithFields(f.errFields)
		}

		return el, nil
	}

	if el.Type().Implements(types.BinaryUnmarshaller) {
		t := el.Interface().(encoding.BinaryUnmarshaler)
		err := t.UnmarshalBinary([]byte(str))
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "unable to unmarshal from text format").
				WithKind(kcderr.Input).
				WithFields(f.errFields)
		}

		return el, nil
	}

	return reflect.Value{}, errors.NewWithKind(kcderr.InputCritical, "an error occur with this getValueFromHttp").
		WithFields(f.errFields)
}

func (f fieldSetter) makeCustomType(str string, ptr bool) (reflect.Value, *errors.Error) {
	el := reflect.New(f.metadata.Type)

	if f.metadata.Type.ConvertibleTo(reflect.TypeOf(time.Duration(0))) {
		duration, err := time.ParseDuration(str)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "unable to parse duration (format: 1ms, 1s, 3h3s)").
				WithKind(kcderr.Input).
				WithFields(f.errFields)
		}

		el.Elem().Set(reflect.ValueOf(duration))

		if ptr {
			return el, nil
		}
		return el.Elem(), nil
	}

	return reflect.Value{}, errors.NewWithKind(kcderr.InputCritical, "an error occur with this getValueFromHttp").
		WithFields(f.errFields)
}
