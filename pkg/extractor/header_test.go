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

func TestHeaderExtractor(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", kcd.Handler(extractorHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	request := e.GET("/")

	addHeader := func(r *httpexpect.Request, assertion extractorAssertion) {
		if reflect.TypeOf(assertion.value).Kind() == reflect.Slice {
			// !!!! Currently header does not support slice values
			return
		}
		r.WithHeader(assertion.rawKey, fmt.Sprintf("%v", assertion.value))
	}

	for _, assertion := range testArray {
		addHeader(request, assertion)
	}

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
