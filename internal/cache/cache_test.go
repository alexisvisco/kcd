package cache

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexisvisco/kcd/internal/types"
)

type structThatImplementTextUnmarshaller struct {
	SuperUser string
}

func (s *structThatImplementTextUnmarshaller) UnmarshalText(text []byte) error {
	s.SuperUser = string(text)
	return nil
}

type structNotUnmarshable struct {
}

type structWithSomething struct {
	User     string `query:"string"`
	UserFunc func() `query:"stringFunc"`
}

type structWithAllSupportedTypes struct {
	String  string  `query:"String"`
	Bool    bool    `query:"Bool"`
	Int     int     `query:"Int"`
	Int8    int8    `query:"Int8"`
	Int16   int16   `query:"Int16"`
	Int32   int32   `query:"Int32"`
	Int64   int64   `query:"Int64"`
	Uint    uint    `query:"Uint"`
	Uint8   uint8   `query:"Uint8"`
	Uint16  uint16  `query:"Uint16"`
	Uint32  uint32  `query:"Uint32"`
	Uint64  uint64  `query:"Uint64"`
	Float32 float32 `query:"Float32"`
	Float64 float64 `query:"Float64"`
}

func TestStructAnalyzer_Cache(t *testing.T) {
	ptrStruct := &structWithSomething{}

	tableTesting := []struct {
		name        string
		st          interface{}
		expectation func(StructCache, *testing.T)
	}{
		{
			"double ptr struct entry (ignored cause double pointer)",
			&ptrStruct,
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 0)
				assert.Len(t, cache.Child, 0)
			},
		},
		{
			"ptr struct entry",
			ptrStruct,
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 1)
				assert.Len(t, cache.Child, 0)
			},
		},
		{
			"string",
			struct {
				Name string `query:"name"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 1)
				assert.Len(t, cache.Child, 0)
				assert.Equal(t, cache.Resolvable[0].Paths["query"], "name")
			},
		},
		{
			"string without tag",
			struct {
				Name string
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 0)
				assert.Len(t, cache.Child, 0)
			},
		},
		{
			"string with default value",
			struct {
				Name string `query:"name" default:"default value"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 1)
				assert.Len(t, cache.Child, 0)
				assert.Equal(t, "name", cache.Resolvable[0].Paths["query"])
				assert.Equal(t, "default value", cache.Resolvable[0].DefaultValue)
			},
		},
		{
			"string with double tag",
			struct {
				Name string `query:"name" path:"namePath"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 1)
				assert.Len(t, cache.Child, 0)
				assert.Equal(t, "name", cache.Resolvable[0].Paths["query"])
				assert.Equal(t, "namePath", cache.Resolvable[0].Paths["path"])
			},
		},
		{
			"pointer string",
			struct {
				Name *string `query:"name"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 1)
				assert.Len(t, cache.Child, 0)
				assert.Equal(t, "name", cache.Resolvable[0].Paths["query"])
			},
		},
		{
			"string array with exploder",
			struct {
				Name []string `query:"name" exploder:","`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 1)
				assert.Len(t, cache.Child, 0)
				assert.Equal(t, "name", cache.Resolvable[0].Paths["query"])
				assert.Equal(t, ",", cache.Resolvable[0].Exploder)
				assert.True(t, cache.Resolvable[0].ArrayOrSlice)
			},
		},
		{
			"simple pointer are accepted while double pointer not",
			struct {
				Name1 []*string  `query:"name1"`
				Name  []**string `query:"name"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 1)
				assert.Len(t, cache.Child, 0)
				assert.Equal(t, "string", cache.Resolvable[0].Type.String())
			},
		},
		{
			"anonymous struct with clear path",
			struct {
				i         int
				Anonymous struct {
					Name string `query:"name"`
				} `query:"anonymous"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 0)
				assert.Len(t, cache.Child, 1)

				childCache := cache.Child[0]
				assert.Len(t, childCache.Resolvable, 1)
				assert.Len(t, childCache.Child, 0)
				assert.Equal(t, "anonymous.name", childCache.Resolvable[0].Paths["query"])
			},
		},
		{
			"anonymous struct with unclear path",
			struct {
				Anonymous struct {
					Name string `query:"name"`
				}
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 0)
				assert.Len(t, cache.Child, 1)

				childCache := cache.Child[0]
				assert.Len(t, childCache.Resolvable, 1)
				assert.Len(t, childCache.Child, 0)
				assert.Equal(t, "name", childCache.Resolvable[0].Paths["query"])
			},
		},
		{
			"struct that implement unmarshaller",
			struct {
				SuperStruct *structThatImplementTextUnmarshaller `query:"superStruct"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 1)
				assert.Len(t, cache.Child, 0)
				assert.Equal(t, "superStruct", cache.Resolvable[0].Paths["query"])
				assert.True(t, cache.Resolvable[0].ImplementUnmarshaller)
			},
		},
		{
			"with a struct that cannot be unmarshalled",
			struct {
				SuperStruct structNotUnmarshable `query:"superStruct"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 0)
				assert.Len(t, cache.Child, 0)
			},
		},
		{
			"struct with all supported types (as embedded)",
			struct {
				structWithAllSupportedTypes
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 0)
				assert.Len(t, cache.Child, 1)
				assert.Len(t, cache.Child[0].Resolvable, len(types.Native))
			},
		},
		{
			"struct with all supported types (as embedded pointer)",
			struct {
				i int
				*structWithAllSupportedTypes
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 0)
				assert.Len(t, cache.Child, 1)
				assert.Equal(t, []int{1}, cache.Child[0].Index)
				assert.Len(t, cache.Child[0].Resolvable, len(types.Native))
			},
		},
		{
			"double pointer ignored",
			struct {
				I **int `query:"super_int"`
			}{},
			func(cache StructCache, t *testing.T) {
				assert.Len(t, cache.Resolvable, 0)
				assert.Len(t, cache.Child, 0)
			},
		},
	}

	for _, assertion := range tableTesting {
		t.Run(assertion.name, func(t *testing.T) {
			analyzer := NewStructAnalyzer([]string{"query", "path", "query"}, []string{}, reflect.TypeOf(assertion.st))
			assertion.expectation(analyzer.Cache(), t)
		})
	}
}
