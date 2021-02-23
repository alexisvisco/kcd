package kcd

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/expectedsh/errors"

	"github.com/expectedsh/kcd/internal/cache"
	"github.com/expectedsh/kcd/internal/decoder"
)

type inputType int

const (
	inputTypeResponse inputType = iota
	inputTypeRequest
	inputTypeInput
	inputTypeCtx
)

var ErrStopHandler = errors.New("KCD_STOP_HANDLER")

// Handler returns a default http handler.
//
// The handler may use the following signature:
//
//  func([response http.ResponseWriter], [request *http.Request], [INPUT object ptr]) ([OUTPUT object], error)
//
// INPUT and OUTPUT struct are both optional.
// As such, the minimal accepted signature is:
//
//  func() error
//
// A complete example for an INPUT struct:
//  type CreateCustomerInput struct {
//		Name          string            `path:"name"`                 // /some-path/{name}
//		Authorization string            `header:"X-authorization"`    // header name 'X-authorization'
//		Emails        []string          `query:"emails" exploder:","` // /some-path/{name}?emails=a@1.fr,b@1.fr
//      Body          map[string]string `json:"body"`                 // json body with {body: {a: "hey", b: "hoy"}}
//
//		ContextualID *struct {
//			ID string `ctx:"id" default:"robot"` // ctx value with key 'id' or it will default set ID to "robot"
//		}
//	}
//
// The wrapped handler will bind the parameters from the query-string,
// path, body and headers, context, and handle the errors.
//
// Handler will panic if the kcd handler or its input/output values
// are of incompatible type.
func Handler(h interface{}, defaultStatusCode int) http.HandlerFunc {
	hv := reflect.ValueOf(h)

	if hv.Kind() != reflect.Func {
		panic(fmt.Sprintf("handler parameters must be a function, got %T", h))
	}
	ht := hv.Type()

	funcName := runtime.FuncForPC(hv.Pointer()).Name()

	orderInput, in := input(ht, funcName)
	out := output(ht, funcName)

	cacheStruct := cache.NewStructAnalyzer(Config.stringsTags(), Config.valuesTags(), in).Cache()

	var input reflect.Value

	// Wrap http handler.
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		// kcd handler has custom input, handle binding.

		if in != nil {
			inputStruct := reflect.New(in)
			input = inputStruct

			// Bind body
			if err := Config.BindHook(w, r, input.Interface()); err != nil {
				Config.ErrorHook(w, r, err, Config.LogHook)
				return
			}

			err := decoder.NewDecoder(r, w, Config.StringsExtractors, Config.ValueExtractors).
				Decode(cacheStruct, input)

			if err != nil {
				Config.ErrorHook(w, r, err, Config.LogHook)
				return
			}

			if err := Config.ValidateHook(r.Context(), inputStruct.Interface()); err != nil {
				Config.ErrorHook(w, r, err, Config.LogHook)
				return
			}
		}

		var err, outputStruct interface{}

		// funcIn contains the input parameters of the kcd handler call.
		var args []reflect.Value
		for _, t := range orderInput {
			switch t {
			case inputTypeInput:
				args = append(args, input)
			case inputTypeRequest:
				args = append(args, reflect.ValueOf(r))
			case inputTypeResponse:
				args = append(args, reflect.ValueOf(w))
			case inputTypeCtx:
				args = append(args, reflect.ValueOf(r.Context()))
			}
		}

		ret := hv.Call(args)
		if out != nil {
			outputStruct = ret[0].Interface()
			err = ret[1].Interface()
		} else {
			err = ret[0].Interface()
		}

		// the handler must stop because its a special error
		if err == ErrStopHandler {
			return
		}

		// Handle the error returned by the handler invocation, if any.
		if err != nil {
			Config.ErrorHook(w, r, err.(error), Config.LogHook)
			return
		}

		// Render the response.
		if err := Config.RenderHook(w, r, outputStruct, defaultStatusCode); err != nil {
			Config.ErrorHook(w, r, err, Config.LogHook)
			return
		}
	}

	return httpHandler
}

var interfaceResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
var interfaceCtx = reflect.TypeOf((*context.Context)(nil)).Elem()

// input checks the input parameters of a kcd handler
// and return the type of the second parameter, if any.
func input(ht reflect.Type, name string) (orderedInputType []inputType, reflectType reflect.Type) {
	n := ht.NumIn()

	if n > 4 {
		panic(fmt.Sprintf(
			"incorrect number of input parameters for handler %s, expected 0 to 4, got %d",
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
		case currentInput.Implements(interfaceCtx):
			if _, exist := setInputType[inputTypeCtx]; exist {
				panic(fmt.Sprintf(
					"invalid parameter %d at handler %s: there is already a context.Context parameter",
					i, name,
				))
			}

			setInputType[inputTypeCtx] = true
			orderedInputType = append(orderedInputType, inputTypeCtx)
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
