package hook

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/expectedsh/errors"
	"github.com/go-chi/chi/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/expectedsh/kcd/internal/kcderr"
)

// ErrorResponse is the default response that send the default error hook.
type ErrorResponse struct {
	ErrorDescription string      `json:"error_description"`
	Error            errors.Kind `json:"error"`

	Fields   map[string]string      `json:"fields,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Error is the default error hook.
// It check the error and return the corresponding response to the client.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("ERROR: ", err) // todo: remove it

	response := ErrorResponse{
		ErrorDescription: "internal server error",
		Error:            errors.KindInternal,
		Fields:           map[string]string{},
		Metadata:         map[string]interface{}{},
	}

	reqID := middleware.GetReqID(r.Context())
	if reqID != "" {
		response.Metadata["request_id"] = reqID
	}

	switch e := err.(type) {
	case validation.Errors:
		w.WriteHeader(http.StatusBadRequest)
		response.Error = errors.KindInvalidArgument
		response.ErrorDescription = "the request has one or multiple invalid fields"

		for k, v := range e {
			response.Fields[k] = v.Error()
		}
	case *errors.Error:
		w.WriteHeader(e.Kind.ToStatusCode())

		response.ErrorDescription = e.Message
		response.Error = e.Kind
	case *kcderr.InputError:

		w.WriteHeader(http.StatusBadRequest)
		response.Error = errors.KindInvalidArgument
		response.ErrorDescription = http.StatusText(http.StatusBadRequest)

		switch e.Extractor {
		case "query", "path", "header", "ctx":
			response.Fields[e.Field] = e.Message
		case "json":
			response.ErrorDescription = e.Message
		}
	case *kcderr.OutputError:
		w.WriteHeader(http.StatusInternalServerError)

		response.Error = errors.KindInternal
		response.ErrorDescription = e.Error()
	}

	// todo: use a log hook to log kcd real (critic) error

	marshal, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Header().Set("Content-type", "application/json")
	_, _ = w.Write(marshal)
}
