package hook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/expectedsh/kcd/pkg/kcderr"
)

// BindHook is the hook called by the wrapped http handler when
// binding an incoming request to the kcd handler's input object.
type Bind func(w http.ResponseWriter, r *http.Request, in interface{}) error

// DefaultBindingHook returns a Bind hook with the default logic, with configurable MaxBodyBytes.
func DefaultBindingHook(maxBodyBytes int64) Bind {
	return func(w http.ResponseWriter, r *http.Request, in interface{}) error {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		if r.ContentLength == 0 {
			return nil
		}

		bytesBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return kcderr.Input{Extractor: "json", Message: "unable to read body"}
		}

		if err := json.Unmarshal(bytesBody, in); err != nil {
			return kcderr.Input{Extractor: "json", Message: "unable to unmarshal request"}
		}

		return nil
	}
}
