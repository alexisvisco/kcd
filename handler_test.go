package kcd

import (
	"math"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"
)

const (
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

type queryExtractorTest struct {
	Str     string  `query:"str"`
	Bool    bool    `query:"bool"`
	Int     int     `query:"int"`
	Int8    int8    `query:"int8"`
	Int16   int16   `query:"int16"`
	Int32   int32   `query:"int32"`
	Int64   int64   `query:"int64"`
	Uint    uint    `query:"uint"`
	Uint8   uint8   `query:"uint8"`
	Uint16  uint16  `query:"uint16"`
	Uint32  uint32  `query:"uint32"`
	Uint64  uint64  `query:"uint64"`
	Float32 float32 `query:"float32"`
	Float64 float64 `query:"float64"`

	ArrStr     []string  `query:"arr_str"`
	ArrBool    []bool    `query:"arr_bool"`
	ArrInt     []int     `query:"arr_int"`
	ArrInt8    []int8    `query:"arr_int8"`
	ArrInt16   []int16   `query:"arr_int16"`
	ArrInt32   []int32   `query:"arr_int32"`
	ArrInt64   []int64   `query:"arr_int64"`
	ArrUint    []uint    `query:"arr_uint"`
	ArrUint8   []uint8   `query:"arr_uint8"`
	ArrUint16  []uint16  `query:"arr_uint16"`
	ArrUint32  []uint32  `query:"arr_uint32"`
	ArrUint64  []uint64  `query:"arr_uint64"`
	ArrFloat32 []float32 `query:"arr_float32"`
	ArrFloat64 []float64 `query:"arr_float64"`

	embeddedQueryExtractor
}

type embeddedQueryExtractor struct {
	StrE string `query:"str_e"`
}

func queryExtractorHandler(req *queryExtractorTest) (response queryExtractorTest, err error) {
	return *req, nil
}

func TestQueryExtractor(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", Handler(queryExtractorHandler, 200))

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("query extractor simple type", func(t *testing.T) {
		expect := e.GET("/").
			WithQuery("str", "string").
			WithQuery("bool", true).
			WithQuery("int", ValInt).
			WithQuery("int8", ValInt8).
			WithQuery("int16", ValInt16).
			WithQuery("int32", ValInt32).
			WithQuery("int64", ValInt64).
			WithQuery("uint", ValUint).
			WithQuery("uint8", ValUint8).
			WithQuery("uint16", ValUint16).
			WithQuery("uint32", ValUint32).
			WithQuery("uint64", ValUint64).
			WithQuery("float32", ValFloat32).
			WithQuery("float64", ValFloat64).Expect()

		j := expect.JSON()
		{
			j.Path("$.Str").Equal("string")
			j.Path("$.Bool").Equal(true)
			j.Path("$.Int").Equal(ValInt)
			j.Path("$.Int8").Equal(ValInt8)
			j.Path("$.Int16").Equal(ValInt16)
			j.Path("$.Int32").Equal(ValInt32)
			j.Path("$.Int64").Equal(ValInt64)
			j.Path("$.Uint").Equal(ValUint)
			j.Path("$.Uint8").Equal(ValUint8)
			j.Path("$.Uint16").Equal(ValUint16)
			j.Path("$.Uint32").Equal(ValUint32)
			j.Path("$.Uint64").Equal(ValUint64)
			j.Path("$.Float32").Equal(ValFloat32)
			j.Path("$.Float64").Equal(ValFloat64)
		}
	})

	t.Run("query extractor array type", func(t *testing.T) {
		expect := e.GET("/").
			WithQuery("arr_str", "string").
			WithQuery("arr_str", "string").
			WithQuery("arr_bool", true).
			WithQuery("arr_bool", true).
			WithQuery("arr_int", ValInt).
			WithQuery("arr_int", ValInt).
			WithQuery("arr_int8", ValInt8).
			WithQuery("arr_int8", ValInt8).
			WithQuery("arr_int16", ValInt16).
			WithQuery("arr_int16", ValInt16).
			WithQuery("arr_int32", ValInt32).
			WithQuery("arr_int32", ValInt32).
			WithQuery("arr_int64", ValInt64).
			WithQuery("arr_int64", ValInt64).
			WithQuery("arr_uint", ValUint).
			WithQuery("arr_uint", ValUint).
			WithQuery("arr_uint8", ValUint8).
			WithQuery("arr_uint8", ValUint8).
			WithQuery("arr_uint16", ValUint16).
			WithQuery("arr_uint16", ValUint16).
			WithQuery("arr_uint32", ValUint32).
			WithQuery("arr_uint32", ValUint32).
			WithQuery("arr_uint64", ValUint64).
			WithQuery("arr_uint64", ValUint64).
			WithQuery("arr_float32", ValFloat32).
			WithQuery("arr_float32", ValFloat32).
			WithQuery("arr_float64", ValFloat64).
			WithQuery("arr_float64", ValFloat64).Expect()

		j := expect.JSON()
		{
			j.Path("$.ArrStr").Equal([]string{"string", "string"})
			j.Path("$.ArrBool").Equal([]bool{true, true})
			j.Path("$.ArrInt").Equal([]int{ValInt, ValInt})
			j.Path("$.ArrInt8").Equal([]int8{ValInt8, ValInt8})
			j.Path("$.ArrInt16").Equal([]int16{ValInt16, ValInt16})
			j.Path("$.ArrInt32").Equal([]int32{ValInt32, ValInt32})
			j.Path("$.ArrInt64").Equal([]int64{ValInt64, ValInt64})
			j.Path("$.ArrUint").Equal([]uint{ValUint, ValUint})
			j.Path("$.ArrUint8").Equal([]uint8{ValUint8, ValUint8})
			j.Path("$.ArrUint16").Equal([]uint16{ValUint16, ValUint16})
			j.Path("$.ArrUint32").Equal([]uint32{ValUint32, ValUint32})
			j.Path("$.ArrUint64").Equal([]uint64{ValUint64, ValUint64})
			j.Path("$.ArrFloat32").Equal([]float32{ValFloat32, ValFloat32})
			j.Path("$.ArrFloat64").Equal([]float64{ValFloat64, ValFloat64})
		}
	})

}
