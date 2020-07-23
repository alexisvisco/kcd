package kcd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/expectedsh/errors"
	"github.com/go-chi/chi/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// defaultValidateHook is the default validation hook.
// It use 'ozzo-validation' to validate structure.
// A structure must implement 'ValidatableWithContext' or 'Validatable'
func defaultValidateHook(ctx context.Context, input interface{}) error {
	var err validation.Errors

	switch v := input.(type) {
	case validation.ValidatableWithContext:
		err = v.ValidateWithContext(ctx).(validation.Errors)
	case validation.Validatable:
		err = v.Validate().(validation.Errors)
	}

	if len(err) == 0 {
		return nil
	}

	return errors.
		NewWithKind(errors.KindInvalidArgument, "the request has one or multiple invalid fields").
		WithMetadata("kcd.fields", err)
}

// defaultRenderHook is the default render hook.
// It marshals the payload to JSON, or returns an empty body if the payload is nil.
func defaultRenderHook(w http.ResponseWriter, _ *http.Request, statusCode int, response interface{}) error {
	if response != nil {
		marshal, err := json.Marshal(response)
		if err != nil {
			return outputError{Err: err}
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

// errorResponse is the default response that send the default error hook
type errorResponse struct {
	ErrorDescription string      `json:"error_description"`
	Error            errors.Kind `json:"error"`

	Fields   map[string]string      `json:"fields,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// defaultErrorHook is the default error hook.
// It check the error and return the corresponding response to the client.
func defaultErrorHook(w http.ResponseWriter, r *http.Request, err error) {

	fmt.Println(err)

	response := errorResponse{
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
	case *errors.Error:
		w.WriteHeader(e.Kind.ToStatusCode())

		response.ErrorDescription = e.Message
		response.Error = e.Kind

		// todo: don't use string literal for kcd.*

		metadata, ok := e.GetMetadata(errorKeyFields)
		if ok {
			m, okMap := metadata.(validation.Errors)
			if okMap {
				for k, v := range m {
					response.Fields[k] = v.Error()
				}
			}
		}

		metadata, ok = e.GetMetadata(errorKeyMetadata)
		if ok {
			m, okMap := metadata.(map[string]interface{})
			if okMap {
				for k, v := range m {
					response.Metadata[k] = v
				}
			}
		}

	case *inputError:

		w.WriteHeader(http.StatusBadRequest)
		response.Error = errors.KindInvalidArgument

		switch e.extractor {
		case queryTag, pathTag, headerTag:
			response.ErrorDescription = http.StatusText(http.StatusBadRequest)
			response.Fields[e.field] = fmt.Sprintf("with %s parameter: %s", e.fieldType, e.message)
		case jsonTag:
			response.ErrorDescription = e.message
		}
	case *outputError:
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

	return
}

// defaultBindHook returns a Bind hook with the default logic, with configurable MaxBodyBytes.
func defaultBindHook(maxBodyBytes int64) BindHook {
	return func(w http.ResponseWriter, r *http.Request, in interface{}) error {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		if r.ContentLength == 0 {
			return nil
		}

		bytesBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return inputError{extractor: jsonTag, message: "unable to read body"}
		}

		if err := json.Unmarshal(bytesBody, in); err != nil {
			return inputError{extractor: jsonTag, message: "unable to unmarshal request"}
		}

		return nil
	}
}
