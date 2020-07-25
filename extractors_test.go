package kcd

import (
	"fmt"
	"math"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"
)

const (
	ValString  = "heyllo"
	ValBool    = true
	ValInt     = int(math.MaxInt32)
	ValInt8    = int8(math.MaxInt8)
	ValInt16   = int16(math.MaxInt16)
	ValInt32   = int32(math.MaxInt32)
	ValInt64   = int64(math.MaxInt64)
	ValUint    = uint(math.MaxUint32)
	ValUint8   = uint8(math.MaxUint8)
	ValUint16  = uint16(math.MaxUint16)
	ValUint32  = uint32(math.MaxUint32)
	ValUint64  = uint64(math.MaxUint64)
	ValFloat32 = float32(1.345876)
	ValFloat64 = float64(1456378934487635.3)
)

type ExtractorTestStruct struct {
	Str     string  `query:"str" path:"str" header:"str"`
	Bool    bool    `query:"bool" path:"bool" header:"bool"`
	Int     int     `query:"int" path:"int" header:"int"`
	Int8    int8    `query:"int8" path:"int8" header:"int8"`
	Int16   int16   `query:"int16" path:"int16" header:"int16"`
	Int32   int32   `query:"int32" path:"int32" header:"int32"`
	Int64   int64   `query:"int64" path:"int64" header:"int64"`
	Uint    uint    `query:"uint" path:"uint" header:"uint"`
	Uint8   uint8   `query:"uint8" path:"uint8" header:"uint8"`
	Uint16  uint16  `query:"uint16" path:"uint16" header:"uint16"`
	Uint32  uint32  `query:"uint32" path:"uint32" header:"uint32"`
	Uint64  uint64  `query:"uint64" path:"uint64" header:"uint64"`
	Float32 float32 `query:"float32" path:"float32" header:"float32"`
	Float64 float64 `query:"float64" path:"float64" header:"float64"`

	PtrStr     *string  `query:"ptr_str" path:"ptr_str" header:"ptr_str"`
	PtrBool    *bool    `query:"ptr_bool" path:"ptr_bool" header:"ptr_bool"`
	PtrInt     *int     `query:"ptr_int" path:"ptr_int" header:"ptr_int"`
	PtrInt8    *int8    `query:"ptr_int8" path:"ptr_int8" header:"ptr_int8"`
	PtrInt16   *int16   `query:"ptr_int16" path:"ptr_int16" header:"ptr_int16"`
	PtrInt32   *int32   `query:"ptr_int32" path:"ptr_int32" header:"ptr_int32"`
	PtrInt64   *int64   `query:"ptr_int64" path:"ptr_int64" header:"ptr_int64"`
	PtrUint    *uint    `query:"ptr_uint" path:"ptr_uint" header:"ptr_uint"`
	PtrUint8   *uint8   `query:"ptr_uint8" path:"ptr_uint8" header:"ptr_uint8"`
	PtrUint16  *uint16  `query:"ptr_uint16" path:"ptr_uint16" header:"ptr_uint16"`
	PtrUint32  *uint32  `query:"ptr_uint32" path:"ptr_uint32" header:"ptr_uint32"`
	PtrUint64  *uint64  `query:"ptr_uint64" path:"ptr_uint64" header:"ptr_uint64"`
	PtrFloat32 *float32 `query:"ptr_float32" path:"ptr_float32" header:"ptr_float32"`
	PtrFloat64 *float64 `query:"ptr_float64" path:"ptr_float64" header:"ptr_float64"`

	ArrStr     []string  `query:"arr_str" path:"arr_str" header:"arr_str"`
	ArrBool    []bool    `query:"arr_bool" path:"arr_bool" header:"arr_bool"`
	ArrInt     []int     `query:"arr_int" path:"arr_int" header:"arr_int"`
	ArrInt8    []int8    `query:"arr_int8" path:"arr_int8" header:"arr_int8"`
	ArrInt16   []int16   `query:"arr_int16" path:"arr_int16" header:"arr_int16"`
	ArrInt32   []int32   `query:"arr_int32" path:"arr_int32" header:"arr_int32"`
	ArrInt64   []int64   `query:"arr_int64" path:"arr_int64" header:"arr_int64"`
	ArrUint    []uint    `query:"arr_uint" path:"arr_uint" header:"arr_uint"`
	ArrUint8   []uint8   `query:"arr_uint8" path:"arr_uint8" header:"arr_uint8"`
	ArrUint16  []uint16  `query:"arr_uint16" path:"arr_uint16" header:"arr_uint16"`
	ArrUint32  []uint32  `query:"arr_uint32" path:"arr_uint32" header:"arr_uint32"`
	ArrUint64  []uint64  `query:"arr_uint64" path:"arr_uint64" header:"arr_uint64"`
	ArrFloat32 []float32 `query:"arr_float32" path:"arr_float32" header:"arr_float32"`
	ArrFloat64 []float64 `query:"arr_float64" path:"arr_float64" header:"arr_float64"`

	EmbeddedExtractorTest
	*EmbeddedPtrExtractorTest
	//embeddedQueryExtractor will don't work since it's not accessible

	SubStruct SubStructExtractorTest

	// will don't work since we need to provide a value to the pointer to inspect its fields.
	// To keep nil = no value, we are not providing this feature.
	// SubStruct *SubStructExtractorTest
}

type EmbeddedExtractorTest struct {
	StrE string `query:"str_embedded" path:"str_embedded" header:"str_embedded"`
}

type EmbeddedPtrExtractorTest struct {
	StrEP string `query:"str_ptr_embedded" path:"str_ptr_embedded" header:"str_ptr_embedded"`
}

type SubStructExtractorTest struct {
	StrS string `query:"sub_struct_str_s" path:"sub_struct_str_s" header:"sub_struct_str_s"`
}

func extractorHandler(req *ExtractorTestStruct) (response ExtractorTestStruct, err error) {
	return *req, nil
}

type extractorAssertion struct {
	rawKey   string
	value    interface{}
	jsonPath string
}

var testArray = []extractorAssertion{
	{"str", ValString, "$.Str"},
	{"bool", ValBool, "$.Bool"},
	{"int", ValInt, "$.Int"},
	{"int8", ValInt8, "$.Int8"},
	{"int16", ValInt16, "$.Int16"},
	{"int32", ValInt32, "$.Int32"},
	{"int64", ValInt64, "$.Int64"},
	{"uint", ValUint, "$.Uint"},
	{"uint8", ValUint8, "$.Uint8"},
	{"uint16", ValUint16, "$.Uint16"},
	{"uint32", ValUint32, "$.Uint32"},
	{"uint64", ValUint64, "$.Uint64"},
	{"float32", ValFloat32, "$.Float32"},
	{"float64", ValFloat64, "$.Float64"},

	{"ptr_str", ValString, "$.Str"},
	{"ptr_bool", ValBool, "$.Bool"},
	{"ptr_int", ValInt, "$.Int"},
	{"ptr_int8", ValInt8, "$.Int8"},
	{"ptr_int16", ValInt16, "$.Int16"},
	{"ptr_int32", ValInt32, "$.Int32"},
	{"ptr_int64", ValInt64, "$.Int64"},
	{"ptr_uint", ValUint, "$.Uint"},
	{"ptr_uint8", ValUint8, "$.Uint8"},
	{"ptr_uint16", ValUint16, "$.Uint16"},
	{"ptr_uint32", ValUint32, "$.Uint32"},
	{"ptr_uint64", ValUint64, "$.Uint64"},
	{"ptr_float32", ValFloat32, "$.Float32"},
	{"ptr_float64", ValFloat64, "$.Float64"},

	{"arr_str", []string{ValString, ValString}, "$.ArrStr"},
	{"arr_bool", []bool{ValBool, ValBool}, "$.ArrBool"},
	{"arr_int", []int{ValInt, ValInt}, "$.ArrInt"},
	{"arr_int8", []int8{ValInt8, ValInt8}, "$.ArrInt8"},
	{"arr_int16", []int16{ValInt16, ValInt16}, "$.ArrInt16"},
	{"arr_int32", []int32{ValInt32, ValInt32}, "$.ArrInt32"},
	{"arr_int64", []int64{ValInt64, ValInt64}, "$.ArrInt64"},
	{"arr_uint", []uint{ValUint, ValUint}, "$.ArrUint"},
	{"arr_uint8", []uint8{ValUint8, ValUint8}, "$.ArrUint8"},
	{"arr_uint16", []uint16{ValUint16, ValUint16}, "$.ArrUint16"},
	{"arr_uint32", []uint32{ValUint32, ValUint32}, "$.ArrUint32"},
	{"arr_uint64", []uint64{ValUint64, ValUint64}, "$.ArrUint64"},
	{"arr_float32", []float32{ValFloat32, ValFloat32}, "$.ArrFloat32"},
	{"arr_float64", []float64{ValFloat64, ValFloat64}, "$.ArrFloat64"},

	{"str_embedded", ValString, "$.StrE"},
	{"str_ptr_embedded", ValString, "$.StrEP"},
	{"sub_struct_str_s", ValString, "$.SubStruct.StrS"},
}

func TestQueryExtractor(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", Handler(extractorHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	request := e.GET("/")

	addQueryParameter := func(r *httpexpect.Request, assertion extractorAssertion) {
		switch reflect.TypeOf(assertion.value).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(assertion.value)

			for i := 0; i < s.Len(); i++ {
				r.WithQuery(assertion.rawKey, s.Index(i))
			}
		default:
			r.WithQuery(assertion.rawKey, assertion.value)
		}
	}

	for _, assertion := range testArray {
		addQueryParameter(request, assertion)
	}

	jsonExpect := request.Expect().JSON()

	for _, assertion := range testArray {
		t.Run(assertion.rawKey, func(t *testing.T) {
			jsonExpect.Path(assertion.jsonPath).Equal(assertion.value)
		})
	}
}

func TestHeaderExtractor(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", Handler(extractorHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	request := e.GET("/")

	addHeader := func(r *httpexpect.Request, assertion extractorAssertion) {
		if reflect.TypeOf(assertion.value).Kind() == reflect.Slice {
			// !!!! Currently header does not support slice values
			return
		}
		r.WithHeader(assertion.rawKey, fmt.Sprintf("%v", assertion.value))
	}

	for _, assertion := range testArray {
		addHeader(request, assertion)
	}

	jsonExpect := request.Expect().JSON()

	for _, assertion := range testArray {
		if reflect.TypeOf(assertion.value).Kind() == reflect.Slice {
			// !!!! Currently header does not support slice values
			continue
		}

		t.Run(assertion.rawKey, func(t *testing.T) {
			jsonExpect.Path(assertion.jsonPath).Equal(assertion.value)
		})
	}
}

func TestPathExtractor(t *testing.T) {
	r := chi.NewRouter()

	urlChi := ""
	urlRequest := ""
	for _, assertion := range testArray {
		if reflect.TypeOf(assertion.value).Kind() == reflect.Slice {
			// !!!! Currently header does not support slice values
			continue
		}
		urlChi += fmt.Sprintf("/{%s}", assertion.rawKey)
		urlRequest += fmt.Sprintf("/%v", assertion.value)
	}

	r.Get(urlChi, Handler(extractorHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	request := e.GET(urlRequest)

	jsonExpect := request.Expect().JSON()

	for _, assertion := range testArray {
		if reflect.TypeOf(assertion.value).Kind() == reflect.Slice {
			// !!!! Currently header does not support slice values
			continue
		}

		t.Run(assertion.rawKey, func(t *testing.T) {
			jsonExpect.Path(assertion.jsonPath).Equal(assertion.value)
		})
	}
}
