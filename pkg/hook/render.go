package hook

import (
	"encoding/json"
	"net/http"

	"github.com/expectedsh/kcd/pkg/kcderr"
)

// Render is the last hook called by the wrapped handler before returning.
// It takes the response, request, the success HTTP status code and the response
// payload as parameters.
//
// Its role is to render the payload to the client to the proper format.
type Render func(w http.ResponseWriter, r *http.Request, defaultSuccessStatusCode int, response interface{}) error

// DefaultRenderHook is the default render hook.
// It marshals the payload to JSON, or returns an empty body if the payload is nil.
func DefaultRenderHook(w http.ResponseWriter, _ *http.Request, statusCode int, response interface{}) error {
	if response != nil {
		marshal, err := json.Marshal(response)
		if err != nil {
			return kcderr.Output{Err: err}
		}

		w.WriteHeader(statusCode)
		if _, err := w.Write(marshal); err != nil {
			return err
		}

	} else {
		w.WriteHeader(statusCode)
	}

	return nil
}
