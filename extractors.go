package kcd

import (
	"net/http"

	"github.com/go-chi/chi"
)

// defaultHeaderExtractor is an extractor that operates on the headers
// of a request.
func defaultHeaderExtractor(_ http.ResponseWriter, r *http.Request, tag string) ([]string, error) {
	header := r.Header.Get(tag)

	if header == "" {
		return nil, nil
	}

	return []string{header}, nil
}

// defaultQueryExtractor is an extractor that operates on the path
// parameters of a request.
func defaultQueryExtractor(_ http.ResponseWriter, r *http.Request, tag string) ([]string, error) {
	return r.URL.Query()[tag], nil
}

// defaultPathExtractor is an extractor that operates on the path
// parameters of a request.
func defaultPathExtractor(_ http.ResponseWriter, r *http.Request, tag string) ([]string, error) {
	p := chi.URLParam(r, tag)
	if p == "" {
		return nil, nil
	}

	return []string{p}, nil
}
