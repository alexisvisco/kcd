package extractor

import (
	"net/http"
)

// DefaultHeaderExtractor is an extractor that operates on the headers
// of a request.
func DefaultHeaderExtractor(_ http.ResponseWriter, r *http.Request, tag string) ([]string, error) {
	header := r.Header.Get(tag)

	return []string{header}, nil
}
