package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/expectedsh/kcd/pkg/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// kcd.Config.BindingHook = ...

	r.Post("/{path}", kcd.Handler(CreateCustomer, http.StatusOK))

	_ = http.ListenAndServe(":3000", r)
}

type CreateCustomerInput struct {
	Name   string   `json:"name"`
	Emails []string `json:"emails"`
}

func (c CreateCustomerInput) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&c.Emails, validation.Each(is.Email)),
	)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request, in *CreateCustomerInput) error {
	return nil
}
