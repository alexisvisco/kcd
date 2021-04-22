package hook_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"

	"github.com/alexisvisco/kcd"
)

type hookRenderStruct struct {
	Name string `json:"name"`
}

func hookRenderEmptyResponseHandler() error {
	return nil
}

func hookRenderPtrHandler(req *hookRenderStruct) (*hookRenderStruct, error) {
	return req, nil
}

func hookRenderHandler(req *hookRenderStruct) (hookRenderStruct, error) {
	return *req, nil
}

func TestRender(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/empty", kcd.Handler(hookRenderEmptyResponseHandler, 200))
	r.Post("/ptr", kcd.Handler(hookRenderPtrHandler, 202))
	r.Post("/", kcd.Handler(hookRenderHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("it should succeed", func(t *testing.T) {
		e.POST("/").WithJSON(hookBindStruct{Name: ValString}).Expect().
			Status(http.StatusOK).
			JSON().Path("$.name").Equal(ValString)
	})

	t.Run("it should succeed with ptr", func(t *testing.T) {
		e.POST("/ptr").WithJSON(hookBindStruct{Name: ValString}).Expect().
			Status(202).
			JSON().Path("$.name").Equal(ValString)
	})

	t.Run("it should succeed with an empty body", func(t *testing.T) {
		e.POST("/empty").Expect().
			Status(http.StatusOK).Body().Equal("")
	})
}
