package kcd

import (
	"context"
	"net/http"
)

const (
	queryTag   = "query"
	pathTag    = "path"
	headerTag  = "header"
	jsonTag    = "json"
	defaultTag = "default"
)

const (
	errorKeyFields   = "kcd.fields"
	errorKeyMetadata = "kcd.metadata"
)

// ErrorHook hook lets you interpret errors returned by your handlers.
// After analysis, the hook should return a suitable http status code
// and and error payload.
// This lets you deeply inspect custom error types.
type ErrorHook func(w http.ResponseWriter, r *http.Request, err error)

// RenderHook is the last hook called by the wrapped handler before returning.
// It takes the response, request, the success HTTP status code and the response
// payload as parameters.
//
// Its role is to render the payload to the client to the proper format.
type RenderHook func(w http.ResponseWriter, r *http.Request, defaultSuccessStatusCode int, response interface{}) error

// BindHook is the hook called by the wrapped http handler when
// binding an incoming request to the kcd handler's input object.
type BindHook func(w http.ResponseWriter, r *http.Request, in interface{}) error

// ValidateHook is the hook called to validate the input.
type ValidateHook func(ctx context.Context, input interface{}) error

// Extractor extract data from request/response with a tag value
type Extractor func(w http.ResponseWriter, r *http.Request, tag string) ([]string, error)

var DefaultExtractors = map[string]Extractor{
	queryTag:  defaultQueryExtractor,
	pathTag:   defaultPathExtractor,
	headerTag: defaultHeaderExtractor,
}

var DefaultHooks = struct {
	Error    ErrorHook
	Render   RenderHook
	Binding  func(maxBodySize int64) BindHook
	Validate ValidateHook
}{
	Error:    defaultErrorHook,
	Render:   defaultRenderHook,
	Binding:  defaultBindHook,
	Validate: defaultValidateHook,
}

var Config = struct {
	Extractors   map[string]Extractor
	ErrorHook    ErrorHook
	RenderHook   RenderHook
	BindHook     BindHook
	ValidateHook ValidateHook
}{

	Extractors: DefaultExtractors,

	ErrorHook:    DefaultHooks.Error,
	RenderHook:   DefaultHooks.Render,
	BindHook:     DefaultHooks.Binding(256 * 1024),
	ValidateHook: DefaultHooks.Validate,
}
