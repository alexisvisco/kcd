package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/expectedsh/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Get("/{name}", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

// CreateCustomerInput is an example of input for an http request.
type CreateCustomerInput struct {
	Name   string   `path:"name"`
	Emails []string `query:"emails" exploder:","`
}

// CustomerOutput is the output type of your handler it contain the input for simplicity.
type CustomerOutput struct {
	Name string `json:"name"`
}

// YourHttpHandler is here using the ctx of the request.
func YourHttpHandler(ctx context.Context, _ *CreateCustomerInput) (CustomerOutput, error) {
	// get the id of the request from the context
	id := middleware.GetReqID(ctx)

	return CustomerOutput{Name: id}, nil
}
