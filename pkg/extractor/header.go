package extractor

import "net/http"

// Header allows to obtain a value from the header of the request.
type Header struct{}

// Extract header from the http request.
func (h Header) Extract(req *http.Request, _ http.ResponseWriter, valueOfTag string) ([]string, error) {
	header := req.Header.Get(valueOfTag)

	if header == "" {
		return nil, nil
	}

	return []string{header}, nil
}

// Tag return the tag name of this extractor.
func (h Header) Tag() string {
	return "header"
}
