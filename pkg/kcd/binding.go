package kcd

import (
	"encoding"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/expectedsh/kcd/pkg/extractor"
	"github.com/expectedsh/kcd/pkg/kcderr"
)

// bind binds the fields the fields of the input object in with
// the values of the parameters extracted from the Gin context.
// It reads tag to know what to extract using the extractor func.
func bind(w http.ResponseWriter, r *http.Request, v reflect.Value, tag string, extract extractor.Extractor) error {
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		field := v.Field(i)

		// Handle embedded fields with a recursive call.
		// If the field is a pointer, but is nil, we
		// create a new value of the same type, or we
		// take the existing memory address.
		if ft.Anonymous {
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
			} else {
				if field.CanAddr() {
					field = field.Addr()
				}
			}

			err := bind(w, r, field, tag, extract)
			if err != nil {
				return err
			}

			continue
		}

		tagValue := ft.Tag.Get(tag)
		if tagValue == "" {
			continue
		}

		bindingError := &kcderr.Input{Field: tagValue, Type: t, Extractor: tag}

		fieldValues, err := extract(w, r, tagValue)
		if err != nil {
			return bindingError.
				WithErr(err).
				WithMessage("unable to extract value from request")
		}

		// Extract default value and use it in place
		// if no values were returned.
		def, ok := ft.Tag.Lookup(Config.DefaultTag)
		if ok && len(fieldValues) == 0 {
			fieldValues = append(fieldValues, def)
		}
		if len(fieldValues) == 0 {
			continue
		}

		// If the field is a nil pointer to a concrete type,
		// create a new addressable value for this type.
		if field.Kind() == reflect.Ptr && field.IsNil() {
			f := reflect.New(field.Type().Elem())
			field.Set(f)
		}

		// Dereference pointer.
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		kind := field.Kind()

		// Multiple values can only be filled to types Slice and Array.
		if len(fieldValues) > 1 && (kind != reflect.Slice && kind != reflect.Array) {
			return bindingError.WithMessage("multiple values not supported")
		}

		// Ensure that the number of values to fill does not exceed the length of a field of type Array.
		if kind == reflect.Array {
			if field.Len() != len(fieldValues) {
				msg := fmt.Sprintf("parameter expect %d values, got %d", field.Len(), len(fieldValues))
				return bindingError.WithMessage(msg)
			}
		}

		if kind == reflect.Slice || kind == reflect.Array {
			// Create a new slice with an adequate
			// length to set all the values.
			if kind == reflect.Slice {
				field.Set(reflect.MakeSlice(field.Type(), 0, len(fieldValues)))
			}
			for i, val := range fieldValues {
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
			continue
		}

		// Fill string value into input field.
		err = bindStringValue(fieldValues[0], field)
		if err != nil {
			return bindingError.
				WithErr(err).
				WithMessage(fmt.Sprintf("unable to set the value %q as type %+v", fieldValues[0], field.Type().Name()))
		}
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
