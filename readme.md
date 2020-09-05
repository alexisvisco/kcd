<p align="center">
	<img width="460" height="300" src="./.github/golang-ss.gif">
</p>
<p align="center">
	<a href="https://github.com/expectedsh/kcd/actions">
		<img width="93" height="20" src="https://github.com/expectedsh/kcd/workflows/Go/badge.svg"></a>
	<a href="https://goreportcard.com/report/github.com/expectedsh/kcd">
		<img width="78" height="20" src="https://goreportcard.com/badge/github.com/expectedsh/kcd"></a>
	<a href='https://coveralls.io/github/expectedsh/kcd'>
		<img src='https://coveralls.io/repos/github/expectedsh/kcd/badge.svg' alt='Coverage Status' /></a>
</p>

------

## :stars: KCD 

KCD is a grandiose REST helper that wrap your shiny handler into a classic http handler. 
It manages all you want for building REST services.

This library is **opinionated** by default but **customizable** which mean it uses some other libraries like Chi, Logrus...
KCD is modular so each pieces of the code that rely on a specific library can be changed. 

## :rocket: QuickStart

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/expectedsh/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	// You can configure kcd with kcd.Config

	r.Get("/{name}", kcd.Handler(YourHttpHandler, http.StatusOK))
	//                       ^ Here the magic happen this is the only thing you need
	//                         to do. Adding kcd.Handler(your handler)
	_ = http.ListenAndServe(":3000", r)
}

// CreateCustomerInput is an example of input for an http request.
type CreateCustomerInput struct {
	Name    string   `path:"name"`                 // you can extract value from: 'path', 'query', 'header', 'ctx'
	Emails  []string `query:"emails" exploder:","` // exploder split value with the characters specified
    Subject string   `json:"body"`                 // it also works with json body
}

// CustomerOutput is the output type of the http request.
type CreateCustomerOutput struct {
	Name string `json:"name"`
}

// YourHttpHandler is your http handler but in a shiny version.
// You can add *http.ResponseWriter or http.Request in params if you want.
func YourHttpHandler(in *CreateCustomerInput) (CreateCustomerOutput, error) {
	// do some stuff here
	fmt.Printf("%+v", in)

	return CreateCustomerOutput{Name: in.Name}, nil
}
```

You can test this code by using curl `curl localhost:3000/supername?emails=alexis@gmail.com,remi@gmail.com`

## :heavy_check_mark: Validation

KCD can validate your input by using a fork of [ozzo-validation](https://github.com/expectedsh/ozzo-validation).

Your input need to implement Validatable or ValidatableWithContext.

```go
// Validate is the function that will be called before calling your shiny handler.
func (c CreateCustomerInput) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&c.Emails, validation.Each(is.Email)),
	)
}
```

## :x: Error handling

KCD handle these kinds of errors: parsing input, validating input and custom handler error.

It uses our internal error package: [errors](https://github.com/expectedsh/errors)

Example of error with a validation failure:

```json
{
  "error_description": "the request has one or multiple invalid fields",
  "error": "invalid_argument",
  "fields": {
    "name": "the length must be between 5 and 20"
  },
  "metadata": {
    "request_id": "BJSIWMOS4o-000001"
  }
}
```

## :coffee: Benefits

- More readable code
- Focus on what it matters: business code
- No more code duplication with unmarshalling, verifying, validating, marshalling ...
- You could have one interface for the client and server implementation

## :muscle: Example

- [*examples/simple/main.go*](./examples/simple/main.go)

