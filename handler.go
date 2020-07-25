package kcd

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

type inputType int

const (
	inputTypeResponse inputType = iota
	inputTypeRequest
	inputTypeInput
)

// Handler returns a default http handler.
//
// The handler may use the following signature:
//
//  func([w http.ResponseWriter], [r *http.Request], [input object ptr]) ([output object], error)
//
// inputError and output objects are both optional.
// As such, the minimal accepted signature is:
//
//  func(w http.ResponseWriter, r *http.Request) error
//
// The wrapped handler will bind the parameters from the query-string,
// path, body and headers, and handle the errors.
//
// Handler will panic if the kcd handler or its input/output values
// are of incompatible type.
func Handler(h interface{}, defaultSuccessStatusCode int) http.HandlerFunc {
	hv := reflect.ValueOf(h)

	if hv.Kind() != reflect.Func {
		panic(fmt.Sprintf("handler parameters must be a httpHandler, got %T", h))
	}
	ht := hv.Type()

	fundName := runtime.FuncForPC(hv.Pointer()).Name()

	orderInput, in := input(ht, fundName)
	out := output(ht, fundName)

	var input *reflect.Value = nil

	// Wrap http handler.
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		// kcd handler has custom input, handle binding.
		if in != nil {
			i := reflect.New(in)
			input = &i

			// Bind body
			if err := Config.BindHook(w, r, input.Interface()); err != nil {
				Config.ErrorHook(w, r, err)
				return
			}

			// Bind query-parameters.
			if err := bind(w, r, i, queryTag, Config.QueryExtractor); err != nil {
				Config.ErrorHook(w, r, err)
				return
			}

			// Bind path arguments.
			if err := bind(w, r, i, pathTag, Config.PathExtractor); err != nil {
				Config.ErrorHook(w, r, err)
				return
			}

			// Bind headers.
			if err := bind(w, r, i, headerTag, Config.HeaderExtractor); err != nil {
				Config.ErrorHook(w, r, err)
				return
			}

			if err := Config.ValidateHook(r.Context(), i.Interface()); err != nil {
				Config.ErrorHook(w, r, err)
				return
			}
		}

		var err, val interface{}

		// funcIn contains the input parameters of the kcd handler call.
		var args []reflect.Value
		for _, t := range orderInput {
			switch t {
			case inputTypeInput:
				args = append(args, *input)
			case inputTypeRequest:
				args = append(args, reflect.ValueOf(r))
			case inputTypeResponse:
				args = append(args, reflect.ValueOf(w))
			}
		}

		ret := hv.Call(args)
		if out != nil {
			val = ret[0].Interface()
			err = ret[1].Interface()
		} else {
			err = ret[0].Interface()
		}

		// Handle the error returned by the handler invocation, if any.
		if err != nil {
			Config.ErrorHook(w, r, err.(error))
			return
		}

		if err := Config.RenderHook(w, r, defaultSuccessStatusCode, val); err != nil {
			Config.ErrorHook(w, r, err)
			return
		}
	}

	return httpHandler
}

var interfaceResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()

// input checks the input parameters of a kcd handler
// and return the type of the second parameter, if any.
func input(ht reflect.Type, name string) (orderedInputType []inputType, reflectType reflect.Type) {
	n := ht.NumIn()

	if n > 3 {
		panic(fmt.Sprintf(
			"incorrect number of input parameters for handler %s, expected 0 or 3, got %d",
			name, n,
		))
	}

	orderedInputType = make([]inputType, 0)
	setInputType := map[inputType]bool{}

	for i := 0; i < n; i++ {
		currentInput := ht.In(i)

		switch {
		case currentInput.Implements(interfaceResponseWriter):
			if _, exist := setInputType[inputTypeResponse]; exist {
				panic(fmt.Sprintf(
					"invalid parameter %d at handler %s: there is already a http.ResponseWriter parameter",
					i, name,
				))
			}

			setInputType[inputTypeResponse] = true
			orderedInputType = append(orderedInputType, inputTypeResponse)
		case currentInput.ConvertibleTo(reflect.TypeOf(&http.Request{})):
			if _, exist := setInputType[inputTypeRequest]; exist {
				panic(fmt.Sprintf(
					"invalid parameter %d at handler %s: there is already a http.Request parameter",
					i, name,
				))
			}

			setInputType[inputTypeRequest] = true
			orderedInputType = append(orderedInputType, inputTypeRequest)
		default:
			if _, exist := setInputType[inputTypeInput]; exist {
				panic(fmt.Sprintf(
					"invalid parameter %d at handler %s: there is already the input parameter",
					i, name,
				))
			}

			if currentInput.Kind() != reflect.Ptr || currentInput.Elem().Kind() != reflect.Struct {
				panic(fmt.Sprintf(
					"invalid %d parameter for handler %s, expected pointer to struct, got %v",
					n, name, currentInput,
				))
			}
			setInputType[inputTypeInput] = true
			orderedInputType = append(orderedInputType, inputTypeInput)
			reflectType = currentInput.Elem()
		}
	}

	return orderedInputType, reflectType
}

// output checks the output parameters of a kcd handler
// and return the type of the return type, if any.
func output(ht reflect.Type, name string) reflect.Type {
	n := ht.NumOut()

	if n < 1 || n > 2 {
		panic(fmt.Sprintf(
			"incorrect number of output parameters for handler %s, expected 1 or 2, got %d",
			name, n,
		))
	}
	// Check the type of the error parameter, which
	// should always come last.
	if !ht.Out(n - 1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic(fmt.Sprintf(
			"unsupported type for handler %s output parameter: expected error interface, got %v",
			name, ht.Out(n-1),
		))
	}
	if n == 2 {
		t := ht.Out(0)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return t
	}
	return nil
}
