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
	Number int      `query:"number"`
	Names  []string `query:"name"`
}

func YourHttpHandler(in *CreateCustomerInput) (string, error) {
	fmt.Printf("%+v", in)

	return fmt.Sprintf("%v ; %v", in.Number, in.Names), nil
}

// Test it : curl -XPOST 'localhost:3000/?number=3&name=vincent&name=remi'
