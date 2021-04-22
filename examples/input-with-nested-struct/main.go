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

type NestedStruct struct {
	Name string `query:"name"`
}

type CustomInput struct {
	Nested NestedStruct `query:"nested"` // ?nested.name=kcd

	NestedStruct // ?name=kcd

	Anonymous struct {
		Key string `query:"key"` // ?key=kcd
	}

	AnonymousValue struct {
		Value        string `query:"key"` // ?anonymous.key=kcd
		NestedStruct        // ?anonymous.name=kcd
	} `query:"anonymous"`
}

func YourHttpHandler(in *CustomInput) (*CustomInput, error) {
	fmt.Printf("%+v", in)

	return in, nil
}

// Test it : curl -XPOST 'localhost:3000?nested.name=alexis&name=remi&key=antoine&anonymous.key=superman&anonymous.name=batman'
