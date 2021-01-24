package main

import (
	"net/http"

	validation "github.com/expectedsh/ozzo-validation/v4"
	"github.com/expectedsh/ozzo-validation/v4/is"
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

func (c *CreateCustomerInput) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, is.Alpha))
}

// CustomerOutput is the output type of your handler it contain the input for simplicity.
type CustomerOutput struct {
	Name string `json:"name"`
}

// YourHttpHandler is here using the ctx of the request.
func YourHttpHandler(c *CreateCustomerInput) (CustomerOutput, error) {
	return CustomerOutput{Name: c.Name}, nil
}
