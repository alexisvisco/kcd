package kcd

import (
	"net/http"
	"strings"

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

// defaultPathExtractor is an extractor that operates on the path
// parameters of a request.
func defaultPathExtractor(_ http.ResponseWriter, r *http.Request, tag string) ([]string, error) {
	p := chi.URLParam(r, tag)
	if p == "" {
		return nil, nil
	}

	return []string{p}, nil
}
