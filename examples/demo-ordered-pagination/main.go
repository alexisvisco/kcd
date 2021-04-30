package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/alexisvisco/kcd/pkg/errors"
	validation "github.com/alexisvisco/ozzo-validation/v4"
	"github.com/go-chi/chi"

	"github.com/alexisvisco/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Post("/{resourceId}", kcd.Handler(YourHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

type Pagination struct {
	Page    uint       `query:"page" default:"1"`
	Limit   uint       `query:"limit" default:"50"`
	OrderBy []*Ordered `query:"orderBy" exploder:"," default:"name:ASC,truc:DESC"`
}

func (p *Pagination) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.OrderBy),
	)
}

type Ordered struct {
	Field string
	Type  string
}

func (o *Ordered) UnmarshalText(text []byte) error {
	list := bytes.Split(text, []byte(":"))
	if len(list) != 2 {
		return errors.NewWithKind(errors.KindInvalidArgument, "order should be of form: field_name:(ASC|DESC)")
	}

	o.Field = string(list[0])
	o.Type = string(list[1])

	return nil
}

func (o *Ordered) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Field,
			validation.Required,
			validation.Length(1, 50),
			validation.In("created_at", "price", "name").Error("Must be created_at, price or name"),
		),
		validation.Field(&o.Type,
			validation.Required,
			validation.In("ASC", "DESC").Error("Must be ASC or DESC"),
		),
	)
}

type ResourceInput struct {
	*Pagination
	ID string `path:"resourceId"`
}

func YourHttpHandler(input *ResourceInput) (*ResourceInput, error) {
	fmt.Printf("%+v", input)
	return input, nil
}

// Test it : curl -XPOST 'localhost:3000'
