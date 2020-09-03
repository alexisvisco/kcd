package extractor

import (
	"net/http"
)

// Strings extract multiples strings values from request/response.
type Strings interface {
	Extract(req *http.Request, res http.ResponseWriter, valueOfTag string) ([]string, error)
	Tag() string
}

// Value extract one value (a type) from http request/response.
type Value interface {
	Extract(req *http.Request, res http.ResponseWriter, valueOfTag string) (interface{}, error)
	Tag() string
}
