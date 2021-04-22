package kcd_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/alexisvisco/kcd"
)

const (
	ValString = "name"
)

type hookBindStruct struct {
	Uint          uint    `path:"uint"`
	JSONName      string  `json:"name"`
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

	StructWithTextUnmarshaller *StructWithTextUnmarshaller `query:"struct_with_text_unmarshaller"`

	SliceStructWithTextUnmarshaller []*StructWithTextUnmarshaller `query:"slice_struct_with_text_unmarshaller"`

	Duration time.Duration `query:"duration"`

	SliceDuration []time.Duration `query:"slice_duration"`

	SlicePointer *[]string `query:"slice_pointer"`

	ArrayDuration [2]time.Duration `query:"array_duration"`

	ArrayInt *[3]int `query:"array_int"`

	StructWithJSONUnmarshaller *StructWithJSONUnmarshaller `query:"struct_with_json_unmarshaller"`

	StructWithBinaryUnmarshaller *StructWithBinaryUnmarshaller `query:"struct_with_binary_unmarshaller"`

	SlicePtrString []*string `query:"slice_ptr_string"`

	SliceInt []int `query:"slice_int"`

	SliceWithUnmarshal *SliceWithUnmarshal `query:"slice_with_unmarshal"`
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

type StructWithTextUnmarshaller struct {
	Value string
}

func (s *StructWithTextUnmarshaller) UnmarshalText(text []byte) error {
	s.Value = strings.ToUpper(string(text))

	return nil
}

func hookBindHandler(bindStruct *hookBindStruct) (hookBindStruct, error) {
	return *bindStruct, nil
}

type StructWithJSONUnmarshaller struct {
	Value string
}

func (s *StructWithJSONUnmarshaller) UnmarshalJSON(bytes []byte) error {
	s.Value = string(bytes)

	return nil
}

type StructWithBinaryUnmarshaller struct {
	Value string
}

func (s *StructWithBinaryUnmarshaller) UnmarshalBinary(data []byte) error {
	s.Value = string(data)
	return nil
}

type SliceWithUnmarshal [2]byte

func (s *SliceWithUnmarshal) UnmarshalText(text []byte) error {
	if len(text) >= 2 {
		s[0] = text[0]
		s[1] = text[1]
	} else {
		s[0] = 'a'
		s[1] = 'b'
	}

	return nil
}

func TestBind(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/{uint}", kcd.Handler(hookBindHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("it should succeed", func(t *testing.T) {
		expect := e.POST("/4").
			WithJSON(hookBindStruct{JSONName: ValString}).
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
			WithQuery("struct_with_text_unmarshaller", "struct_with_text_unmarshaller").
			WithQuery("slice_struct_with_text_unmarshaller", "one").
			WithQuery("slice_struct_with_text_unmarshaller", "two").
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
			WithQuery("struct_with_json_unmarshaller", `{"hi": 1}`).
			WithQuery("struct_with_binary_unmarshaller", `goijerierjoer`).
			WithQuery("slice_ptr_string", `hey`).
			WithQuery("slice_ptr_string", `how`).
			WithQuery("slice_int", 0).
			WithQuery("slice_with_unmarshal", "hg").
			Expect()

		raw := expect.Body().Raw()

		assert.NotContains(t, raw, "EmbeddedStringPtrNil")

		j := expect.JSON()
		j.Path("$.Uint").Equal(4)
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
		j.Path("$.StructWithTextUnmarshaller.Value").Equal("STRUCT_WITH_TEXT_UNMARSHALLER")
		j.Path("$.SliceStructWithTextUnmarshaller").Equal([]map[string]interface{}{
			{"Value": "ONE"},
			{"Value": "TWO"},
		})
		j.Path("$.Duration").Equal(time.Second * 45)
		j.Path("$.SliceDuration").Equal([]time.Duration{time.Second * 45, time.Second * 2})
		j.Path("$.SlicePointer").Equal([]string{"one", "two"})
		j.Path("$.ArrayDuration").Equal([]time.Duration{time.Second * 45, time.Second * 2})
		j.Path("$.ArrayInt").Equal([]int{1, 2, 3})
		j.Path("$.StructWithJSONUnmarshaller.Value").Equal(`{"hi": 1}`)
		j.Path("$.StructWithBinaryUnmarshaller.Value").Equal(`goijerierjoer`)
		j.Path("$.SlicePtrString").Equal([]string{"hey", "how"})
		j.Path("$.SliceInt").Equal([]int{0})
		j.Path("$.SliceWithUnmarshal").Equal(SliceWithUnmarshal{'h', 'g'})
	})

	t.Run("it should fail because of invalid body", func(t *testing.T) {
		e.POST("/3").
			WithHeader("Content-type", "application/json").
			WithBytes([]byte("{ this is a malformed json")).Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.error_description").Equal("unable to read json request")
	})

	t.Run("it should succeed with an empty body", func(t *testing.T) {
		e.POST("/3").Expect().
			Status(http.StatusOK).
			JSON().Path("$.name").Equal("")
	})

	t.Run("it should fail because of invalid uint", func(t *testing.T) {
		e.POST("/-3").Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.fields.uint").Equal("invalid positive integer")
	})

	t.Run("it should fail because of invalid boolean", func(t *testing.T) {
		e.POST("/4").WithQuery("query_bool", "yes").Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.fields.query_bool").Equal("invalid boolean")
	})

	t.Run("it should fail because of invalid float", func(t *testing.T) {
		e.POST("/4").WithQuery("query_float", "yes").Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.fields.query_float").Equal("invalid floating number")
	})

	t.Run("it should fail because of invalid duration", func(t *testing.T) {
		e.POST("/4").WithQuery("duration", "yes").Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.fields.duration").Equal("unable to parse duration (format: 1ms, 1s, 3h3s)")
	})

	t.Run("it should fail because of invalid duration", func(t *testing.T) {
		e.POST("/4").WithQuery("slice_int", "yes").Expect().
			Status(http.StatusBadRequest).
			JSON().Path("$.fields.slice_int").Equal("invalid integer")
	})
}
