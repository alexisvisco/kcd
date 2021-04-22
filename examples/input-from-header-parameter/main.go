package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/alexisvisco/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Post("/", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

type CreateCustomerInput struct {
	Number float64 `header:"x-float-id"`
}

func YourHttpHandler(in *CreateCustomerInput) (float64, error) {
	fmt.Printf("%+v", in)

	return in.Number, nil
}

// Test it : curl -XPOST 'localhost:3000' -H "x-float-id: 3.5"
