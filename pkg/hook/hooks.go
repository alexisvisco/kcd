package hook

import (
	"context"
	"net/http"
)

// ErrorHook hook lets you interpret errors returned by your handlers.
// After analysis, the hook should return a suitable http status code
// and and error payload.
// This lets you deeply inspect custom error types.
type ErrorHook func(w http.ResponseWriter, r *http.Request, err error, logger LogHook)

// RenderHook is the last hook called by the wrapped handler before returning.
// It takes the response, request, the success HTTP status code and the response
// payload as parameters.
//
// Its role is to render the payload to the client to the proper format.
type RenderHook func(w http.ResponseWriter, r *http.Request, response interface{}, defaultStatusCode int) error

// BindHook is the hook called by the wrapped http handler when
// binding an incoming request to the kcd handler's input object.
type BindHook func(w http.ResponseWriter, r *http.Request, in interface{}) error

// ValidateHook is the hook called to validate the input.
// The default expected return (handled by the error hook) is a map[string]error.
type ValidateHook func(ctx context.Context, input interface{}) error

// LogHook is the logger triggered after the error hook.
// It can show detailed error about a problem that you can't explain to users.
type LogHook func(w http.ResponseWriter, r *http.Request, err error)
