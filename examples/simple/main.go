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

	// You can configure kcd with kcd.Config.  ErrorHook,
	//                                         RenderHook,
	//                                         BindHook,
	//                                         ValidateHook,
	//                                         LogHook,
	//                                         StringsExtractors,
	//                                         ValueExtractors

	r.Get("/{name}", kcd.Handler(YourHttpHandler, http.StatusOK))
	//                       ^ Here the magic happen this is the only thing you need
	//                         to do. Adding kcd.Handler(your handler)
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

// YourHttpHandler is your http handler but in a shiny version.
// You can add *http.ResponseWriter or http.Request in params if you want.
func YourHttpHandler(in *CreateCustomerInput) (CustomerOutput, error) {
	// do some stuff here

	fmt.Printf("%+v", in)

	return CustomerOutput{Name: in.Name}, nil
}
