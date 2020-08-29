package extractor_test

import (
	"testing"
)

func TestCtxExtractor(t *testing.T) {
	//r := chi.NewRouter()
	//r.Use(func(handler http.Handler) http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		ctx := r.Context()
	//		for _, assertion := range testArray {
	//			ctx = context.WithValue(ctx, assertion.rawKey, assertion.value)
	//		}
	//		handler.ServeHTTP(w, r.WithContext(ctx))
	//	})
	//})
	//
	//r.Get("/", kcd.Handler(extractorHandler, 200))
	//
	//server := httptest.NewServer(r)
	//defer server.Close()
	//
	//e := httpexpect.New(t, server.URL)
	//
	//request := e.GET("/")
	//
	//ex := request.Expect()
	//fmt.Println(ex.Raw())
	//jsonExpect := ex.JSON()
	//
	//for _, assertion := range testArray {
	//	t.Run(assertion.rawKey, func(t *testing.T) {
	//		jsonExpect.Path(assertion.jsonPath).Equal(assertion.value)
	//	})
	//}
}
