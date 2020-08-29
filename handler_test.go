package kcd_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/expectedsh/kcd"
)

const (
	ValString = "name"
)

type hookBindStruct struct {
	JsonName      string  `json:"name"`
	QueryString   string  `query:"query_string"`
	QueryFloat    float32 `query:"query_float"`
	QueryBool     bool    `query:"query_bool"`
	QueryInt      int     `query:"query_int"`
	QueryIntPtr   *int    `query:"query_int"`
	QueryListUint []uint  `query:"query_list_uint"`
	Default       string  `default:"default_value"`

	AnonymousStruct struct {
		AnonymousField string `header:"header_anonymous_struct"`
	}

	Embedded

	*EmbeddedPtr

	*EmbeddedPtrNil

	FilledStruct *StructOne

	StructWithUnmarshaller *StructWithUnmarshaller `query:"struct_with_unmarshaller"`

	SliceStructWithUnmarshaller []*StructWithUnmarshaller `query:"slice_struct_with_unmarshaller"`

	Duration time.Duration `query:"duration"`

	SliceDuration []time.Duration `query:"slice_duration"`

	SlicePointer *[]string `query:"slice_pointer"`

	ArrayDuration [2]time.Duration `query:"array_duration"`

	ArrayInt *[3]int `query:"array_int"`
}

type Embedded struct {
	EmbeddedString string `query:"embedded_string"`
}

type EmbeddedPtr struct {
	EmbeddedStringPtr string `query:"embedded_string_ptr"`
}

type EmbeddedPtrNil struct {
	EmbeddedStringPtrNil string `query:"embedded_string_ptr_nil"`
}

type StructOne struct {
	Inner *StructTwo
}

type StructTwo struct {
	FilledField string `query:"filled_field"`
}

type StructWithUnmarshaller struct {
	Value string
}

func (s *StructWithUnmarshaller) UnmarshalText(text []byte) error {
	s.Value = strings.ToUpper(string(text))

	return nil
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
		expect := e.POST("/4").
			WithJSON(hookBindStruct{JsonName: ValString}).
			WithQuery("query_string", "query_string").
			WithQuery("query_float", 1.4).
			WithQuery("query_bool", true).
			WithQuery("query_int", 50).
			WithQuery("query_list_uint", 0).
			WithQuery("query_list_uint", 1).
			WithHeader("header_anonymous_struct", "header_anonymous_struct").
			WithQuery("embedded_string", "embedded_string").
			WithQuery("embedded_string_ptr", "embedded_string_ptr").
			WithQuery("filled_field", "filled_field").
			WithQuery("struct_with_unmarshaller", "struct_with_unmarshaller").
			WithQuery("slice_struct_with_unmarshaller", "one").
			WithQuery("slice_struct_with_unmarshaller", "two").
			WithQuery("duration", "45s").
			WithQuery("slice_duration", "45s").
			WithQuery("slice_duration", "2s").
			WithQuery("slice_pointer", "one").
			WithQuery("slice_pointer", "two").
			WithQuery("array_duration", "45s").
			WithQuery("array_duration", "2s").
			WithQuery("array_int", "1").
			WithQuery("array_int", "2").
			WithQuery("array_int", "3").
			Expect()

		raw := expect.Body().Raw()
		fmt.Println(raw)

		assert.NotContains(t, raw, "EmbeddedStringPtrNil")

		j := expect.JSON()
		j.Path("$.name").Equal(ValString)
		j.Path("$.QueryString").Equal("query_string")
		j.Path("$.QueryFloat").Equal(1.4)
		j.Path("$.QueryBool").Equal(true)
		j.Path("$.QueryInt").Equal(50)
		j.Path("$.QueryIntPtr").Equal(50)
		j.Path("$.QueryListUint").Equal([]uint{0, 1})
		j.Path("$.Default").Equal("default_value")
		j.Path("$.AnonymousStruct.AnonymousField").Equal("header_anonymous_struct")
		j.Path("$.EmbeddedString").Equal("embedded_string")
		j.Path("$.EmbeddedStringPtr").Equal("embedded_string_ptr")
		j.Path("$.FilledStruct.Inner.FilledField").Equal("filled_field")
		j.Path("$.StructWithUnmarshaller.Value").Equal("STRUCT_WITH_UNMARSHALLER")
		j.Path("$.SliceStructWithUnmarshaller").Equal([]map[string]interface{}{
			{"Value": "ONE"},
			{"Value": "TWO"},
		})
		j.Path("$.Duration").Equal(time.Second * 45)
		j.Path("$.SliceDuration").Equal([]time.Duration{time.Second * 45, time.Second * 2})
		j.Path("$.SlicePointer").Equal([]string{"one", "two"})
		j.Path("$.ArrayDuration").Equal([]time.Duration{time.Second * 45, time.Second * 2})
		j.Path("$.ArrayInt").Equal([]int{1, 2, 3})

	})

	t.Run("it should fail because of invalid body", func(t *testing.T) {
		e.POST("/3").WithBytes([]byte("{ this is a malformed json")).Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.error_description").Equal("unable to read json request")
	})

	t.Run("it should succeed with an empty body", func(t *testing.T) {
		e.POST("/3").Expect().
			Status(http.StatusOK).
			JSON().Path("$.name").Equal("")
	})
}
