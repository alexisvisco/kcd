package extractor

import (
	"net/http"
	"strings"
)

// DefaultQueryExtractor is an extractor that operates on the path
// parameters of a request.
func DefaultQueryExtractor(_ http.ResponseWriter, r *http.Request, tag string) ([]string, error) {
	var params []string
	query := r.URL.Query()[tag]

	splitFn := func(c rune) bool {
		return c == ','
	}

	for _, q := range query {
		params = append(params, strings.FieldsFunc(q, splitFn)...)
	}

	return params, nil
}
