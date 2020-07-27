package kcd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"

	"github.com/expectedsh/kcd"
)

const (
	ValString = "hey"
)

type hookBindStruct struct {
	Name           string  `json:"name"`
	Query          string  `query:"query"`
	QueryFloat     float32 `query:"qf"`
	Bool           bool    `query:"bool"`
	QueryInt       int     `query:"queryInt"`
	PathUInt       uint    `query:"uint"`
	QueryList      []uint  `query:"qlist"`
	TestDefault    string  `json:"test_default" default:"4car"`
	EmbeddedString struct {
		Test string `header:"test"`
	}

	HookBindE
}

type HookBindE struct {
	String string `query:"hey"`
}

func hookBindHandler(bindStruct *hookBindStruct) (hookBindStruct, error) {
	return *bindStruct, nil
}

func TestBind(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/{uint}", kcd.Handler(hookBindHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("it should succeed", func(t *testing.T) {
		j := e.POST("/4").
			WithQuery("query", "sucre").
			WithQuery("queryInt", "5678").
			WithQuery("qf", "4.4").
			WithQuery("qlist", "1").
			WithQuery("bool", "false").
			WithQuery("qlist", "2").
			WithHeader("header", "sucre2").
			WithQuery("hey", "hook").
			WithJSON(hookBindStruct{Name: ValString}).
			Expect().
			Status(http.StatusOK).
			JSON()

		j.Path("$.name").Equal(ValString)
		j.Path("$.test_default").Equal("4car")
	})

	t.Run("it should fail because of invalid body", func(t *testing.T) {
		e.POST("/3").WithBytes([]byte("{ this is a malformed json")).Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.error_description").Equal("unable to unmarshal request")
	})

	t.Run("it should succeed with an empty body", func(t *testing.T) {
		e.POST("/3").Expect().
			Status(http.StatusOK).
			JSON().Path("$.name").Equal("")
	})
}
