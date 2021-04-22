package hook

import (
	"encoding/json"
	"net/http"

	"github.com/expectedsh/errors"

	"github.com/alexisvisco/kcd/internal/kcderr"
)

// Render is the default render hook.
// It marshals the output into JSON, or returns an empty body if the payload is nil.
func Render(w http.ResponseWriter, _ *http.Request, response interface{}, statusCode int) error {
	if response != nil {
		marshal, err := json.Marshal(response)
		if err != nil {
			return errors.Wrap(err, "unable to render response in json format").WithKind(kcderr.OutputCritical)
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(statusCode)
		if _, err := w.Write(marshal); err != nil {
			return errors.Wrap(err, "unable to write response").WithKind(kcderr.OutputCritical)
		}
	} else {
		w.WriteHeader(statusCode)
	}

	return nil
}
