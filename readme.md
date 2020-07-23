<p align="center">
  <img width="460" height="300" src="./.github/golang-ss.gif">
</p>

### KCD ?

KCD is a grandiose http handler that manages un-marshall, validating, errors, marshaling ... Opinionated by default but fully customizable.

It wraps your shiny handler in a http.HandlerFunc. 

### Example

*examples/simple/main.go:*
```go
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

func CreateCustomer(in *CreateCustomerInput) (CreateCustomerInput, error) {
	return *in, nil
}
```

#### Opinionated

KCD use:
- [github.com/go-chi/chi](https://github.com/go-chi/chi) for the router
- [github.com/expectedsh/errors](https://github.com/expectedsh/errors) for the library errors
- [github.com/go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation) for the validation of structs

#### Customizable

KCD works with a hooks & extractor system that are exposed and editable in the main package via the variable `Config`.

In this Config variable you have theses variables:

| Variable Name | Description |
|---|---|
| QueryExtractor | Used to retrieve query parameter, works with http stdlib |
| PathExtractor | Used to retrieve path variable, chi is used|
| HeaderExtractor | Used to retrieve header value, works with http stdlib |
| ErrorHook | This is the way the REST app will manage error, by default it return JSON error, errors (from expected) is used|
| RenderHook | Use json as response |
| BindHook | Use json tag of request struct retrieved with the body of the request |
| ValidateHook | Use ozzo-validation to validate request struct |


