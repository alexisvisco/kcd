package hook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	validation "github.com/expectedsh/ozzo-validation/v4"
	"github.com/expectedsh/ozzo-validation/v4/is"
	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"

	"github.com/expectedsh/kcd"
)

type hookValidateStruct struct {
	Name string `json:"name"`
}

func (h *hookValidateStruct) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.Name, is.Digit))
}

type hookValidateStructCtx struct {
	Name string `json:"name"`
}

func (h *hookValidateStructCtx) ValidateWithContext(_ context.Context) error {
	return validation.ValidateStruct(h,
		validation.Field(&h.Name, is.Alpha))
}

func hookValidateHandler(req *hookValidateStruct) error {
	return nil
}

func hookValidateHandlerCtx(req *hookValidateStructCtx) error {
	return nil
}

func TestValidate(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/", kcd.Handler(hookValidateHandler, 200))
	r.Post("/ctx", kcd.Handler(hookValidateHandlerCtx, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("it should succeed without error", func(t *testing.T) {
		e.POST("/").WithJSON(hookValidateStruct{Name: "0123"}).Expect().
			Status(http.StatusOK)
	})

	t.Run("it should succeed with ctx and without error", func(t *testing.T) {
		e.POST("/ctx").WithJSON(hookValidateStructCtx{Name: "abc"}).Expect().
			Status(http.StatusOK)
	})

	t.Run("it should fail because validation fail", func(t *testing.T) {
		e.POST("/").WithJSON(hookValidateStruct{Name: "a"}).Expect().
			Status(http.StatusBadRequest).JSON().Path("$.fields.name").Equal(is.ErrDigit.Message())
	})
}
