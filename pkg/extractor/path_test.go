package extractor_test

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"

	"github.com/alexisvisco/kcd"
)

func TestPathExtractor(t *testing.T) {
	r := chi.NewRouter()

	urlChi := ""
	urlRequest := ""
	for _, assertion := range testArray {
		if reflect.TypeOf(assertion.value).Kind() == reflect.Slice {
			// !!!! Currently header does not support slice values
			continue
		}
		urlChi += fmt.Sprintf("/{%s}", assertion.rawKey)
		urlRequest += fmt.Sprintf("/%v", assertion.value)
	}

	r.Get(urlChi, kcd.Handler(extractorHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	request := e.GET(urlRequest)

	jsonExpect := request.Expect().JSON()

	for _, assertion := range testArray {
		if reflect.TypeOf(assertion.value).Kind() == reflect.Slice {
			// !!!! Currently header does not support slice values
			continue
		}

		t.Run(assertion.rawKey, func(t *testing.T) {
			jsonExpect.Path(assertion.jsonPath).Equal(assertion.value)
		})
	}
}
