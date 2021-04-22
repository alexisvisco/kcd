package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/alexisvisco/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Post("/{names}", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

type CreateCustomerInput struct {
	Heroes []string `query:"heroes" exploder:","`
	Names  []string `path:"names" exploder:":"`
}

func YourHttpHandler(in *CreateCustomerInput) (*CreateCustomerInput, error) {
	fmt.Printf("%+v", in)

	return in, nil
}

// Test it : curl -XPOST 'localhost:3000/alexis:remi:antoine?heroes=superman,batman,flash'
