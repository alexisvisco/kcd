package hook_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"

	"github.com/expectedsh/kcd"
)

type hookBindStruct struct {
	Name string `json:"name"`
}

func hookBindHandler(bindStruct *hookBindStruct) (hookBindStruct, error) {
	return *bindStruct, nil
}

func TestBind(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/", kcd.Handler(hookBindHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("it should succeed", func(t *testing.T) {
		e.POST("/").WithJSON(hookBindStruct{Name: ValString}).Expect().
			Status(http.StatusOK).
			JSON().Path("$.name").Equal(ValString)
	})

	t.Run("it should fail because of invalid body", func(t *testing.T) {
		e.POST("/").WithBytes([]byte("{ this is a malformed json")).Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.error_description").Equal("unable to unmarshal request")
	})

	t.Run("it should succeed with an empty body", func(t *testing.T) {
		e.POST("/").Expect().
			Status(http.StatusOK).
			JSON().Path("$.name").Equal("")
	})
}
