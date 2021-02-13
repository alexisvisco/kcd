package main

import (
	"net/http"

	validation "github.com/expectedsh/ozzo-validation/v4"
	"github.com/expectedsh/ozzo-validation/v4/is"
	"github.com/go-chi/chi"

	"github.com/expectedsh/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{name}", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

// CreateCustomerInput is an example of input for an http request.
type CreateCustomerInput struct {
	Name string `path:"name"`
}

func (c *CreateCustomerInput) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, is.Alpha))
}

// YourHttpHandler is here using the ctx of the request.
func YourHttpHandler(c *CreateCustomerInput) (*CreateCustomerInput, error) {
	return c, nil
}

// Test it : curl localhost:3000/alexis
