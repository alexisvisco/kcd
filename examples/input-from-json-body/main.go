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
	Name string `json:"name"`
}

func YourHttpHandler(in *CreateCustomerInput) (string, error) {
	fmt.Printf("%+v", in)

	return in.Name, nil
}

// Test it : curl -XPOST -H "Content-type: application/json" -d '{"name": "alexis"}' 'localhost:3000'
