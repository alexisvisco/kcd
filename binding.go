package kcd

import (
	"encoding"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/expectedsh/kcd/internal/kcderr"
	"github.com/expectedsh/kcd/pkg/extractor"
)

type binder struct {
	response          http.ResponseWriter
	request           *http.Request
	stringsExtractors []extractor.Strings
	valueExtractors   []extractor.Value
}

func newBinder(response http.ResponseWriter, request *http.Request,
	stringsExtractors []extractor.Strings, valueExtractors []extractor.Value) *binder {
	return &binder{
		response:          response,
		request:           request,
		stringsExtractors: stringsExtractors,
		valueExtractors:   valueExtractors,
	}
}

func (b *binder) bind(v reflect.Value) error {
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)  // reflect.StructField
		fieldValue := v.Field(i) // reflect.Value

		if fieldType.Anonymous {
			if err := b.handleEmbeddedField(fieldValue); err != nil {
				return err
			}
			continue
		}

		if fieldValue.Kind() == reflect.Struct {
			err := b.bind(fieldValue)
			if err != nil {
				return err
			}
			continue
		}

		bindingError := &kcderr.InputError{FieldType: t, Field: t.Name()}

		values, err := b.extractValue(fieldType, bindingError)
		if err != nil {
			return err
		}

		if len(values) == 0 {
			defaultValue, ok := fieldType.Tag.Lookup("default")
			if ok {
				values = []string{defaultValue}
			}
		}

		if len(values) == 0 {
			continue
		}

		// If the fieldValue is a nil pointer to a concrete type,
		// create a new addressable value for this type.
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			f := reflect.New(fieldValue.Type().Elem())
			fieldValue.Set(f)
		}

		// Dereference pointer.
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}
		kind := fieldValue.Kind()

		// Multiple values can only be filled to types Slice and Array.
		if len(values) > 1 && (kind != reflect.Slice && kind != reflect.Array) {
			return bindingError.WithMessage("multiple values not supported")
		}

		// Ensure that the number of values to fill does not exceed the length of a fieldValue of type Array.
		if kind == reflect.Array {
			if fieldValue.Len() != len(values) {
				msg := fmt.Sprintf("parameter expect %d values, got %d", fieldValue.Len(), len(values))
				return bindingError.WithMessage(msg)
			}
		}

		if kind == reflect.Slice || kind == reflect.Array {
			if err := b.handleSliceAndArray(fieldValue, values, bindingError); err != nil {
				return err
			}
			continue
		}

		// Fill string value into input fieldValue.
		err = bindStringValue(values[0], fieldValue)
		if err != nil {
			return bindingError.
				WithErr(err).
				WithMessage(fmt.Sprintf("unable to set the value %q as type %+v", values[0], fieldValue.Type().Name()))
		}
	}

	return nil
}

func (b *binder) extractValue(fieldType reflect.StructField, bindingError *kcderr.InputError) ([]string, error) {
	for _, ex := range b.stringsExtractors {
		tag := ex.Tag()
		tagValue := fieldType.Tag.Get(tag)
		if tagValue == "" {
			continue
		}

		bindingError.Field = tagValue
		bindingError.Extractor = tag

		values, err := ex.Extract(b.request, b.response, tagValue)

		if err != nil {
			return nil, bindingError.
				WithErr(err).
				WithMessage("unable to extract value from request")
		}

		if len(values) > 0 {
			return values, nil
		}
	}

	return nil, nil
}

func (b *binder) handleSliceAndArray(field reflect.Value, values []string,
	bindingError *kcderr.InputError) (err error) {
	kind := field.Kind()
	// Create a new slice with an adequate
	// length to set all the values.
	if kind == reflect.Slice {
		field.Set(reflect.MakeSlice(field.Type(), 0, len(values)))
	}
	for i, val := range values {
		v := reflect.New(field.Type().Elem()).Elem()
		err = bindStringValue(val, v)
		if err != nil {
			return bindingError.
				WithErr(err).
				WithMessage(fmt.Sprintf("unable to set the value %q as type %+v", val, v.Type().Name()))
		}
		if kind == reflect.Slice {
			field.Set(reflect.Append(field, v))
		}
		if kind == reflect.Array {
			field.Index(i).Set(v)
		}
	}
	return nil
}

// handleEmbeddedField embedded fields with a call to bind.
// If the field is a pointer, but is nil, we create a new value of the same type, or we
// take the existing memory address.
func (b *binder) handleEmbeddedField(field reflect.Value) (err error) {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
	} else {
		if field.CanAddr() {
			field = field.Addr()
		}
	}

	err = b.bind(field)
	if err != nil {
		return err
	}

	return nil
}

// bindStringValue converts and bind the value s to the the reflected value v.
func bindStringValue(s string, v reflect.Value) error {
	// Ensure that the reflected value is addressable
	// and wasn't obtained by the use of an unexported
	// struct field, or calling a setter will panic.
	if !v.CanSet() {
		return fmt.Errorf("unaddressable value: %v", v)
	}

	i := reflect.New(v.Type()).Interface()

	// If the value implements the encoding.TextUnmarshaller
	// interface, bind the returned string representation.
	if unmarshaller, ok := i.(encoding.TextUnmarshaler); ok {
		if err := unmarshaller.UnmarshalText([]byte(s)); err != nil {
			return err
		}
		v.Set(reflect.Indirect(reflect.ValueOf(unmarshaller)))
		return nil
	}

	// Handle time.Duration.
	if _, ok := i.(time.Duration); ok {
		d, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(d))
	}

	// Switch over the kind of the reflected value
	// and convert the string to the proper type.
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetFloat(i)
	default:
		return fmt.Errorf("unsupported parameter type: %v", v.Kind())
	}
	return nil
}
