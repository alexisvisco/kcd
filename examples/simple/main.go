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

	// kcd.Config.BindHook = ...

	r.Post("/{name}", kcd.Handler(CreateCustomer, http.StatusOK))

	_ = http.ListenAndServe(":3000", r)
}

type CreateCustomerInput struct {
	Name   string   `json:"name" path:"name"`
	Emails []string `json:"emails"`
}

func (c CreateCustomerInput) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&c.Emails, validation.Each(is.Email)),
	)
}

type Customer struct {
	CreateCustomerInput
}

func CreateCustomer(in *CreateCustomerInput) (Customer, error) {
	// do some stuff here

	return Customer{*in}, nil
}
