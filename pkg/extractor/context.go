package extractor

import (
	"net/http"
)

// Context extract value from the the context of the request.
type Context struct{}

// Extract value from the context of the request.
func (c Context) Extract(req *http.Request, _ http.ResponseWriter, valueOfTag string) (interface{}, error) {
	return req.Context().Value(valueOfTag), nil
}

// Tag return the tag name of this extractor.
func (c Context) Tag() string {
	return "ctx"
}
