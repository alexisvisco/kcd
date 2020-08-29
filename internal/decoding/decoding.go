package decoding

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/expectedsh/kcd/internal/cache"
	"github.com/expectedsh/kcd/pkg/extractor"
)

type Decoder struct {
	req *http.Request
	res http.ResponseWriter

	stringsExtractors []extractor.Strings
	valueExtractors   []extractor.Value
}

func NewDecoder(
	req *http.Request,
	res http.ResponseWriter,
	stringsExtractors []extractor.Strings,
	valueExtractors []extractor.Value,
) *Decoder {
	return &Decoder{
		req:               req,
		res:               res,
		stringsExtractors: stringsExtractors,
		valueExtractors:   valueExtractors,
	}
}

type previousFields struct {
	root          reflect.Value
	uninitialized [][]int
}

func (d previousFields) getCurrentValue() reflect.Value {
	var field = d.root

	for _, fieldIndex := range d.uninitialized {

		if field.Kind() == reflect.Ptr {
			field = field.Elem().FieldByIndex(fieldIndex)
		} else {
			field = field.FieldByIndex(fieldIndex)
		}

		if field.Kind() == reflect.Ptr {
			val := reflect.New(field.Type().Elem())
			field.Set(val)

			field = field.Elem()
		}
	}

	return field
}

func (d Decoder) Decode(c cache.StructCache, root reflect.Value) error {
	return d.decode(c, root.Type(), previousFields{root: root})
}

func (d Decoder) decode(c cache.StructCache, root reflect.Type, prev previousFields) error {
	fieldsToSet := make([]setterContext, 0, len(c.Resolvable))

	for _, metadata := range c.Resolvable {
		tag, path, v, err := d.value(metadata)
		if err != nil {
			return err
		}

		if v != nil {

			fieldsToSet = append(fieldsToSet, setterContext{
				tag:      tag,
				path:     path,
				metadata: metadata,
				value:    v,
			})
		}
	}

	if len(fieldsToSet) > 0 {
		currentValue := prev.getCurrentValue()

		prev = previousFields{
			root:          currentValue,
			uninitialized: [][]int{},
		}

		for _, setterCtx := range fieldsToSet {
			var field reflect.Value

			if currentValue.Kind() == reflect.Ptr {
				field = currentValue.Elem().FieldByIndex(setterCtx.metadata.Index)
			} else {
				field = currentValue.FieldByIndex(setterCtx.metadata.Index)
			}

			if err := newFieldSetter(field, setterCtx).set(); err != nil {
				return err
			}
		}
	}

	for _, structCache := range c.Child {
		tmpRoot := root

		newPreviousFields := previousFields{
			root:          prev.root,
			uninitialized: prev.uninitialized,
		}

		if root.Kind() == reflect.Ptr {
			tmpRoot = root.Elem().FieldByIndex(structCache.Index).Type
			newPreviousFields.uninitialized = append(newPreviousFields.uninitialized, structCache.Index)
		} else {
			tmpRoot = root.FieldByIndex(structCache.Index).Type
			newPreviousFields.uninitialized = append(newPreviousFields.uninitialized, structCache.Index)
		}

		if err := d.decode(structCache, tmpRoot, newPreviousFields); err != nil {
			return err
		}
	}

	return nil

}

func (d Decoder) value(r cache.FieldMetadata) (tag, key string, v interface{}, err error) {
	for _, e := range d.stringsExtractors {
		path, ok := r.Paths[e.Tag()]
		if ok {
			list, err := e.Extract(d.req, d.res, path)
			if err != nil {
				return "", "", nil, err
			}

			if len(list) == 0 {
				continue
			}

			if len(r.Exploder) > 0 && len(list) == 1 && r.ArrayOrSlice {
				list = strings.Split(list[0], r.Exploder)
			}

			return e.Tag(), path, list, nil
		}
	}

	for _, e := range d.valueExtractors {
		path, ok := r.Paths[e.Tag()]
		if ok {
			v, err := e.Extract(d.req, d.res, path)
			if err != nil {
				return "", "", nil, err
			}

			if len(r.Exploder) > 0 && r.ArrayOrSlice {
				switch t := v.(type) {
				case string:
					list := strings.Split(t, r.Exploder)
					if len(list) > 1 {
						return "", "", list, nil
					}
				}
			}

			if v != nil {
				return e.Tag(), path, v, nil
			}
		}
	}

	if len(r.DefaultValue) > 0 {
		def := r.DefaultValue

		if len(r.Exploder) > 0 && r.ArrayOrSlice {
			list := strings.Split(def, r.Exploder)
			if len(list) > 1 {
				return "default", "", list, nil
			}
		}

		return "default", "", def, nil
	}

	return "", "", nil, nil
}
