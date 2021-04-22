package extractor_test

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"

	"github.com/alexisvisco/kcd"
)

func TestQueryExtractor(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", kcd.Handler(extractorHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	request := e.GET("/")

	addQueryParameter := func(r *httpexpect.Request, assertion extractorAssertion) {
		switch reflect.TypeOf(assertion.value).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(assertion.value)

			for i := 0; i < s.Len(); i++ {
				r.WithQuery(assertion.rawKey, s.Index(i))
			}
		default:
			r.WithQuery(assertion.rawKey, assertion.value)
		}
	}

	for _, assertion := range testArray {
		addQueryParameter(request, assertion)
	}

	jsonExpect := request.Expect().JSON()

	for _, assertion := range testArray {
		t.Run(assertion.rawKey, func(t *testing.T) {
			jsonExpect.Path(assertion.jsonPath).Equal(assertion.value)
		})
	}
}
