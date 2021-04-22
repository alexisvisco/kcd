package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/alexisvisco/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Post("/{number}", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

type CreateCustomerInput struct {
	Number uint `path:"number"`
}

func YourHttpHandler(in *CreateCustomerInput) (uint, error) {
	fmt.Printf("%+v", in)

	return in.Number, nil
}

// Test it : curl -XPOST 'localhost:3000/1'
