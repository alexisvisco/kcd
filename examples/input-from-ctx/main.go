package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/alexisvisco/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "test", "this is a value from a context")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Post("/", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

type CreateCustomerInput struct {
	RequestID string `ctx:"test"`
}

func YourHttpHandler(in *CreateCustomerInput) (string, error) {
	fmt.Printf("%+v\n", in)

	return in.RequestID, nil
}

// Test it : curl -XPOST 'localhost:3000'
