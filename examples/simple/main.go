package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/expectedsh/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// kcd.Configuration.BindHook = ...

	r.Post("/{name}", kcd.Handler(SuperShinyHandler, http.StatusOK))

	_ = http.ListenAndServe(":3000", r)
}

// CreateCustomerInput is an example of input for an http request.
type CreateCustomerInput struct {
	Name   string   `json:"name" path:"name"`
	Emails []string `json:"emails"`
}

// Validate is the function that will be called before calling your shiny handler.
func (c CreateCustomerInput) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&c.Emails, validation.Each(is.Email)),
	)
}

// Customer is the output type of your handler it contain the input for simplicity.
type Customer struct {
	CreateCustomerInput
}

// SuperShinyHandler is your http handler but in a shiny version.
func SuperShinyHandler(in *CreateCustomerInput) (Customer, error) {
	// do some stuff here

	return Customer{*in}, nil
}
