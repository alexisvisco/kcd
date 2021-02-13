package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/expectedsh/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

type Output struct {
	Name string `json:"name"`
}

func YourHttpHandler() (Output, error) {
	return Output{
		Name: "Hello world",
	}, nil
}

// Test it : curl 'localhost:3000'
