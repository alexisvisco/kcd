package extractor

import (
	"net/http"
)

// An extractorFunc extracts data from a gin context according to
// parameters specified in a field tag.
type Extractor func(w http.ResponseWriter, r *http.Request, str string) ([]string, error)
