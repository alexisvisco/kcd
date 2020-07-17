package hook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/expectedsh/errors"
	"github.com/go-chi/chi/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/expectedsh/kcd/pkg/kcderr"
)

// Error hook lets you interpret errors returned by your handlers.
// After analysis, the hook should return a suitable http status code
// and and error payload.
// This lets you deeply inspect custom error types.
type Error func(w http.ResponseWriter, r *http.Request, err error)

type errorResponse struct {
	ErrorDescription string      `json:"error_description"`
	Error            errors.Kind `json:"error"`

	Fields   map[string]string      `json:"fields,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func DefaultErrorHook(w http.ResponseWriter, r *http.Request, err error) {

	response := errorResponse{
		ErrorDescription: "internal server error",
		Fields:           map[string]string{},
		Metadata:         map[string]interface{}{},
	}

	reqID := middleware.GetReqID(r.Context())
	if reqID != "" {
		response.Metadata["request_id"] = reqID
	}

	switch e := err.(type) {
	case *errors.Error:
		w.WriteHeader(e.Kind.ToStatusCode())

		response.ErrorDescription = e.Message
		response.Error = e.Kind

		// todo: don't use string literal for kcd.*

		metadata, ok := e.GetMetadata("kcd.fields")
		if ok {
			m, okMap := metadata.(validation.Errors)
			if okMap {
				for k, v := range m {
					response.Fields[k] = v.Error()
				}
			}
		}

		metadata, ok = e.GetMetadata("kcd.metadata")
		if ok {
			m, okMap := metadata.(map[string]interface{})
			if okMap {
				for k, v := range m {
					response.Metadata[k] = v
				}
			}
		}

	case kcderr.Input:
		w.WriteHeader(http.StatusBadRequest)
		response.Error = errors.KindInvalidArgument

		// todo: don't use string literal for query, path, header, json

		switch e.Extractor {
		case "query", "path", "header":
			response.ErrorDescription = http.StatusText(http.StatusBadRequest)
			response.Fields[e.Field] = fmt.Sprintf("with %s parameter: %s", e.Type, e.Message)
		case "json":
			response.ErrorDescription = e.Message
		}
	case kcderr.Output:
		w.WriteHeader(http.StatusInternalServerError)

		response.Error = errors.KindInternal
		response.ErrorDescription = e.Error()
	}

	// todo: use a log hook to log kcd real (critic) error

	marshal, err := json.Marshal(response)
	if err != nil {
		return
	}

	_, _ = w.Write(marshal)

	return
}
