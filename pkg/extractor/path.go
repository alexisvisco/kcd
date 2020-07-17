package extractor

import (
	"net/http"

	"github.com/go-chi/chi"
)

// DefaultPathExtractor is an extractor that operates on the path
// parameters of a request.
func DefaultPathExtractor(_ http.ResponseWriter, r *http.Request, tag string) ([]string, error) {
	p := chi.URLParam(r, tag)

	return []string{p}, nil
}
