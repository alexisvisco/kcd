package extractor

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Path extract value from the chi router.
type Path struct{}

// Extract value from the chi router.
func (p Path) Extract(req *http.Request, _ http.ResponseWriter, valueOfTag string) ([]string, error) {
	str := chi.URLParam(req, valueOfTag)
	if str == "" {
		return nil, nil
	}

	return []string{str}, nil
}

// Tag return the tag name of this extractor.
func (p Path) Tag() string {
	return "path"
}
