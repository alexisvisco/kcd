package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/expectedsh/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Post("/", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

type CreateCustomerInput struct {
	Power      float64 `header:"power" default:"7.5"`
	IsSuperman *bool   `query:"superman" default:"true"`
	CanFly     bool    `query:"canFly" default:"false"`
}

func YourHttpHandler(in *CreateCustomerInput) (*CreateCustomerInput, error) {
	fmt.Printf("%+v", in)

	return in, nil
}

// Test it :
//   - curl -XPOST 'localhost:3000'
//     Power, IsSuperman and Can fly will be set with all their default value
//   - curl -XPOST 'localhost:3000?canFly=true'
//     Only Power and IsSuperman will be set with default values
