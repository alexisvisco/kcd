package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/expectedsh/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	// kcd.Config.ErrorHook ...

	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "id", 12345)
			handler.ServeHTTP(w, r)
		})
	})

	r.Get("/{name}", kcd.Handler(SuperShinyHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

// CreateCustomerInput is an example of input for an http request.
type CreateCustomerInput struct {
	Name         string   `path:"name"`
	Emails       []string `query:"emails" exploder:","`
	ContextualID *struct {
		ID int `ctx:"id"`
	}
}

// Validate is the function that will be called before calling your shiny handler.
func (c CreateCustomerInput) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&c.Emails, validation.Each(is.Email)),
		validation.Field(&c.ContextualID, validation.Required),
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
