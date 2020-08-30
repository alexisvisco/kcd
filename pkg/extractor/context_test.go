package extractor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/expectedsh/kcd"
)

type chocolate struct {
	ID   int
	Name string
}

type ctxRequest struct {
	Choco *chocolate `ctx:"choco"`
}

func TestCtxExtractor(t *testing.T) {
	r := chi.NewRouter()
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// nolint
			ctx = context.WithValue(ctx, "choco", &chocolate{
				ID:   123,
				Name: "kcd",
			})
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Get("/", kcd.Handler(func(req *ctxRequest) error {
		assert.NotNil(t, req.Choco)
		assert.Equal(t, 123, req.Choco.ID)
		assert.Equal(t, "kcd", req.Choco.Name)

		return nil
	}, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	httpexpect.New(t, server.URL).GET("/").Expect().Status(200)
}
