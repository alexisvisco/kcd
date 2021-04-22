package hook_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexisvisco/kcd/pkg/errors"
	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/alexisvisco/kcd"
)

type hookErrorStruct struct {
	Value int16 `json:"value" query:"value"`
}

func hookErrorHandler(errorStruct *hookErrorStruct) error {
	if errorStruct.Value < 1 {
		return errors.
			NewWithKind(errors.KindUnavailable, "value is unavailable").
			WithField("value", errorStruct.Value)
	}

	if errorStruct.Value == 50 {
		return fmt.Errorf("test")
	}

	return nil
}

func TestError(t *testing.T) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Post("/x", kcd.Handler(hookErrorHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("it should use query input error", func(t *testing.T) {
		e.POST("/x").WithQuery("value", "ab").Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.fields.value").Equal("invalid integer")
	})

	t.Run("it should use expectedsh/errors", func(t *testing.T) {
		e.POST("/x").Expect().
			Status(http.StatusServiceUnavailable).
			JSON().Path("$.error_description").Equal("value is unavailable")
	})

	t.Run("it should use normal error", func(t *testing.T) {
		e.POST("/x").WithQuery("value", "50").Expect().
			Status(http.StatusInternalServerError)
	})
}
