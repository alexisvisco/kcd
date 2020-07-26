package extractor

import "net/http"

// Query allows to obtain a value from the query params of the request.
type Query struct{}

// Extract query params from the http request.
func (q Query) Extract(req *http.Request, _ http.ResponseWriter, valueOfTag string) ([]string, error) {
	return req.URL.Query()[valueOfTag], nil
}

// Tag return the tag name of this extractor.
func (q Query) Tag() string {
	return "query"
}
