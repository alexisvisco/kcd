package extractor

import (
	"net/http"
)

type Strings interface {
	Extract(req *http.Request, res http.ResponseWriter, valueOfTag string) ([]string, error)
	Tag() string
}

type Value interface {
	Extract(req *http.Request, res http.ResponseWriter, valueOfTag string) (interface{}, error)
	Tag() string
}
