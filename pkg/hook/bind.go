package hook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/expectedsh/kcd/internal/kcderr"
)

// Bind returns a Bind hook, it will read only maxBodyBytes bytes from the body and unmarshall
// the input interface with the json encoding of the stdlib.
func Bind(maxBodyBytes int64) BindHook {
	return func(w http.ResponseWriter, r *http.Request, in interface{}) error {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		if r.ContentLength == 0 {
			return nil
		}

		bytesBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return kcderr.InputError{Extractor: "json", Message: "unable to read body"}
		}

		if err := json.Unmarshal(bytesBody, in); err != nil {
			return kcderr.InputError{Extractor: "json", Message: "unable to unmarshal request"}
		}

		return nil
	}
}
