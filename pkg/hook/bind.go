package hook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/alexisvisco/kcd/pkg/errors"

	"github.com/alexisvisco/kcd/internal/kcderr"
)

// Bind returns a Bind hook, it will read only maxBodyBytes bytes from the body and unmarshall
// the input interface with the json encoding of the stdlib.
func Bind(maxBodyBytes int64) BindHook {
	return func(w http.ResponseWriter, r *http.Request, in interface{}) error {
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			if r.ContentLength == 0 {
				return nil
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

			bytesBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrap(err, "unable to read body").WithKind(kcderr.InputCritical)
			}

			if err := json.Unmarshal(bytesBody, in); err != nil {
				return errors.Wrap(err, "unable to read json request").
					WithKind(kcderr.Input).
					WithField("decoding-strategy", "json")
			}
		}

		return nil
	}
}
