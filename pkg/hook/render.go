package hook

import (
	"encoding/json"
	"net/http"

	"github.com/expectedsh/kcd/internal/kcderr"
)

// Render is the default render hook.
// It marshals the output into JSON, or returns an empty body if the payload is nil.
func Render(w http.ResponseWriter, _ *http.Request, response interface{}, statusCode int) error {
	if response != nil {
		marshal, err := json.Marshal(response)
		if err != nil {
			return kcderr.OutputError{Err: err}
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(statusCode)
		if _, err := w.Write(marshal); err != nil {
			return err
		}
	} else {
		w.WriteHeader(statusCode)
	}

	return nil
}
