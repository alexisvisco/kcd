package hook

import (
	"encoding/json"
	"net/http"

	"github.com/expectedsh/errors"
	validation "github.com/expectedsh/ozzo-validation/v4"
	"github.com/go-chi/chi/middleware"

	"github.com/expectedsh/kcd/internal/kcderr"
)

var ErrorHookStatusCodeMinLogged = 500

// ErrorResponse is the default response that send the default error hook.
type ErrorResponse struct {
	ErrorDescription string      `json:"error_description"`
	Error            errors.Kind `json:"error"`

	Fields   map[string]string      `json:"fields,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Error is the default error hook.
// It check the error and return the corresponding response to the client.
// logger parameter is optional (you can set it to nil)
func Error(w http.ResponseWriter, r *http.Request, err error, logger LogHook) {
	response := ErrorResponse{
		ErrorDescription: "internal server error",
		Error:            errors.KindInternal,
		Fields:           map[string]string{},
		Metadata:         map[string]interface{}{},
	}

	w.Header().Set("Content-type", "application/json")

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
		if e.Kind == kcderr.Input {
			w.WriteHeader(http.StatusBadRequest)
			response.Error = errors.KindInvalidArgument
			response.ErrorDescription = http.StatusText(http.StatusBadRequest)

			// TODO(alexis) 23/08/2020: maybe handle ctx tag as a internal server error because it is handled by the
			// 							input provided by the developer and it is not an user input.

			tag, _ := e.GetField("tag")
			path, _ := e.GetField("path")

			switch tag {
			case "query", "path", "header", "ctx", "default":
				response.Fields[path.(string)] = e.Message
			case "json":
				response.ErrorDescription = e.Message
			}

			break
		}

		if e.Kind == kcderr.InputCritical {
			w.WriteHeader(http.StatusInternalServerError)
			response.Error = e.Kind
			response.ErrorDescription = e.Message

			if logger != nil {
				logger(w, r, e)
			}

			break
		}

		if e.Kind == kcderr.OutputCritical {
			w.WriteHeader(http.StatusInternalServerError)
			response.Error = e.Kind
			response.ErrorDescription = e.Message

			if logger != nil {
				logger(w, r, e)
			}

			break
		}

		w.WriteHeader(e.Kind.ToStatusCode())

		response.ErrorDescription = e.Message
		response.Error = e.Kind

		if e.Kind.ToStatusCode() >= ErrorHookStatusCodeMinLogged {
			// ensure there is an error internal if the status code is ErrorHookStatusCodeMinLogged (for instance when omission of the kind)
			response.Error = errors.KindInternal
		}

		if e.Kind.ToStatusCode() >= 500 {
			if logger != nil {
				logger(w, r, e)
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response.Error = errors.KindInternal

		if logger != nil {
			logger(w, r, e)
		}
	}

	marshal, err := json.Marshal(response)
	if err != nil {
		return
	}

	_, _ = w.Write(marshal)
}
