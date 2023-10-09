/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package binding

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	ginBinder "github.com/gin-gonic/gin/binding"
	"github.com/hertz-contrib/binding/go_tagexpr"
	//"github.com/valyala/fasthttp"
)

type Req struct {
	I int `query:"i"`
	J int `query:"j"`
	K int `query:"k"`
}

//func getBenchCtx() Ctx {
//	app := New()
//	ctx := app.NewCtx(&fasthttp.RequestCtx{}).(*DefaultCtx)
//	var u = fasthttp.URI{}
//	u.SetQueryString("i=1&j=1&k=1")
//	ctx.Request().SetURI(&u)
//
//	return ctx
//}

/***** normal ******/
func Benchmark_HertzNormalQuery(b *testing.B) {
	ctx := app.NewContext(0)
	ctx.URI().SetQueryString("i=1&j=1&k=1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		err := ctx.Bind(&v)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
	}
}

func Benchmark_GoTagexprNormalQuery(b *testing.B) {
	ctx := app.NewContext(0)
	ctx.URI().SetQueryString("i=1&j=1&k=1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		err := go_tagexpr.NewBinder().BindAndValidate(&ctx.Request, &v, nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
	}
}

//func Benchmark_FiberNormalQuery(b *testing.B) {
//	ctx := getBenchCtx()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		var v = Req{}
//		err := ctx.Bind().Req(&v).Err()
//		if err != nil {
//			b.Error(err)
//			b.FailNow()
//		}
//	}
//}

func Benchmark_GinNormalQuery(b *testing.B) {
	type Req struct {
		I int `form:"i"`
		J int `form:"j"`
		K int `form:"k"`
	}
	httpReq, _ := http.NewRequest("GET", "", nil)
	params := make(url.Values)
	params.Add("i", "1")
	params.Add("j", "1")
	params.Add("k", "1")
	httpReq.URL.RawQuery = params.Encode()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		binder := ginBinder.Form
		err := binder.Bind(httpReq, &v)
		if err != nil {
			b.Error("error")
		}
	}
}

/***** BigQuerySmallField ******/
func Benchmark_HertzBigQuerySmallField(b *testing.B) {
	ctx := app.NewContext(0)
	str := strings.Builder{}
	for i := 0; i < 100; i++ {
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
	}
	ctx.URI().SetQueryString(str.String() + "i=1&j=1&k=1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		err := ctx.Bind(&v)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if v.I != 1 {
			b.FailNow()
		}
	}
}

func Benchmark_GoTagexprBigQuerySmallField(b *testing.B) {
	ctx := app.NewContext(0)
	str := strings.Builder{}
	for i := 0; i < 100; i++ {
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
	}
	ctx.URI().SetQueryString(str.String() + "i=1&j=1&k=1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		err := go_tagexpr.NewBinder().BindAndValidate(&ctx.Request, &v, nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if v.I != 1 {
			b.FailNow()
		}
	}
}

//func getBenchCtxBigQuerySmallField() Ctx {
//	app := New()
//	ctx := app.NewCtx(&fasthttp.RequestCtx{}).(*DefaultCtx)
//	var u = fasthttp.URI{}
//	str := strings.Builder{}
//	for i := 0; i < 100; i++ {
//		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
//	}
//	u.SetQueryString(str.String() + "i=1&j=1&k=1")
//	ctx.Request().SetURI(&u)
//
//	return ctx
//}

//func Benchmark_FiberBigQuerySmallField(b *testing.B) {
//	ctx := getBenchCtxBigQuerySmallField()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		var v = Req{}
//		err := ctx.Bind().Req(&v).Err()
//		if err != nil {
//			b.Error(err)
//			b.FailNow()
//		}
//		if v.I != 1 {
//			b.FailNow()
//		}
//	}
//}

func Benchmark_GinBigQuerySmallField(b *testing.B) {
	type Req struct {
		I int `form:"i"`
		J int `form:"j"`
		K int `form:"k"`
	}
	httpReq, _ := http.NewRequest("GET", "", nil)
	params := make(url.Values)
	params.Add("i", "1")
	params.Add("j", "1")
	params.Add("k", "1")
	for i := 0; i < 100; i++ {
		params.Add(fmt.Sprintf("h%d", i), fmt.Sprint(i))
	}
	httpReq.URL.RawQuery = params.Encode()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		binder := ginBinder.Form
		err := binder.Bind(httpReq, &v)
		if err != nil {
			b.Error("error")
		}
		if v.I != 1 {
			b.FailNow()
		}
	}
}

/***** BigQueryBigField ******/
type BigFieldReq struct {
	I   int `query:"i"`
	J   int `query:"j"`
	K   int `query:"k"`
	I0  int `query:"i"`
	I1  int `query:"i"`
	I2  int `query:"i"`
	I3  int `query:"i"`
	I4  int `query:"i"`
	I5  int `query:"i"`
	I6  int `query:"i"`
	I7  int `query:"i"`
	I8  int `query:"i"`
	I9  int `query:"i"`
	I10 int `query:"i"`
	I11 int `query:"i"`
	I12 int `query:"i"`
	I13 int `query:"i"`
	I14 int `query:"i"`
	I15 int `query:"i"`
	I16 int `query:"i"`
	I17 int `query:"i"`
	I18 int `query:"i"`
	I19 int `query:"i"`
	I20 int `query:"i"`
	I21 int `query:"i"`
	I22 int `query:"i"`
	I23 int `query:"i"`
	I24 int `query:"i"`
	I25 int `query:"i"`
	I26 int `query:"i"`
	I27 int `query:"i"`
	I28 int `query:"i"`
	I29 int `query:"i"`
	I30 int `query:"i"`
	I31 int `query:"i"`
	I32 int `query:"i"`
	I33 int `query:"i"`
	I34 int `query:"i"`
	I35 int `query:"i"`
	I36 int `query:"i"`
	I37 int `query:"i"`
	I38 int `query:"i"`
	I39 int `query:"i"`
	I40 int `query:"i"`
	I41 int `query:"i"`
	I42 int `query:"i"`
	I43 int `query:"i"`
	I44 int `query:"i"`
	I45 int `query:"i"`
	I46 int `query:"i"`
	I47 int `query:"i"`
	I48 int `query:"i"`
	I49 int `query:"i"`
	I50 int `query:"i"`
	I51 int `query:"i"`
	I52 int `query:"i"`
	I53 int `query:"i"`
	I54 int `query:"i"`
	I55 int `query:"i"`
	I56 int `query:"i"`
	I57 int `query:"i"`
	I58 int `query:"i"`
	I59 int `query:"i"`
	I60 int `query:"i"`
	I61 int `query:"i"`
	I62 int `query:"i"`
	I63 int `query:"i"`
	I64 int `query:"i"`
	I65 int `query:"i"`
	I66 int `query:"i"`
	I67 int `query:"i"`
	I68 int `query:"i"`
	I69 int `query:"i"`
	I70 int `query:"i"`
	I71 int `query:"i"`
	I72 int `query:"i"`
	I73 int `query:"i"`
	I74 int `query:"i"`
	I75 int `query:"i"`
	I76 int `query:"i"`
	I77 int `query:"i"`
	I78 int `query:"i"`
	I79 int `query:"i"`
	I80 int `query:"i"`
	I81 int `query:"i"`
	I82 int `query:"i"`
	I83 int `query:"i"`
	I84 int `query:"i"`
	I85 int `query:"i"`
	I86 int `query:"i"`
	I87 int `query:"i"`
	I88 int `query:"i"`
	I89 int `query:"i"`
	I90 int `query:"i"`
	I91 int `query:"i"`
	I92 int `query:"i"`
	I93 int `query:"i"`
	I94 int `query:"i"`
	I95 int `query:"i"`
	I96 int `query:"i"`
}

func Test_PrintBigFieldReq(t *testing.T) {
	str := strings.Builder{}
	str.WriteString("type BigFieldReq struct {\n")
	for i := 0; i < 97; i++ {
		str.WriteString(fmt.Sprintf("I%d []int `form:\"i\"`\n", i))
	}
	str.WriteString("}\n")
	fmt.Println(str.String())
}

func Benchmark_HertzBigQueryBigField(b *testing.B) {
	ctx := app.NewContext(0)
	str := strings.Builder{}
	for i := 0; i < 100; i++ {
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
	}
	ctx.URI().SetQueryString(str.String() + "i=1&j=1&k=1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := BigFieldReq{}
		err := ctx.Bind(&v)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if v.I != 1 {
			b.FailNow()
		}
	}
}

func Benchmark_GoTagexprBigQueryBigField(b *testing.B) {
	ctx := app.NewContext(0)
	str := strings.Builder{}
	for i := 0; i < 100; i++ {
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
	}
	ctx.URI().SetQueryString(str.String() + "i=1&j=1&k=1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := BigFieldReq{}
		err := go_tagexpr.NewBinder().BindAndValidate(&ctx.Request, &v, nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if v.I != 1 {
			b.FailNow()
		}
	}
}

//func Benchmark_FiberBigQueryBigField(b *testing.B) {
//	ctx := getBenchCtxBigQuerySmallField()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		var v = BigFieldReq{}
//		err := ctx.Bind().Req(&v).Err()
//		if err != nil {
//			b.Error(err)
//			b.FailNow()
//		}
//		if v.I != 1 {
//			b.FailNow()
//		}
//	}
//}

func Benchmark_GinBigQueryBigField(b *testing.B) {
	type Req struct {
		I   int `form:"i"`
		J   int `form:"j"`
		K   int `form:"k"`
		I0  int `form:"i"`
		I1  int `form:"i"`
		I2  int `form:"i"`
		I3  int `form:"i"`
		I4  int `form:"i"`
		I5  int `form:"i"`
		I6  int `form:"i"`
		I7  int `form:"i"`
		I8  int `form:"i"`
		I9  int `form:"i"`
		I10 int `form:"i"`
		I11 int `form:"i"`
		I12 int `form:"i"`
		I13 int `form:"i"`
		I14 int `form:"i"`
		I15 int `form:"i"`
		I16 int `form:"i"`
		I17 int `form:"i"`
		I18 int `form:"i"`
		I19 int `form:"i"`
		I20 int `form:"i"`
		I21 int `form:"i"`
		I22 int `form:"i"`
		I23 int `form:"i"`
		I24 int `form:"i"`
		I25 int `form:"i"`
		I26 int `form:"i"`
		I27 int `form:"i"`
		I28 int `form:"i"`
		I29 int `form:"i"`
		I30 int `form:"i"`
		I31 int `form:"i"`
		I32 int `form:"i"`
		I33 int `form:"i"`
		I34 int `form:"i"`
		I35 int `form:"i"`
		I36 int `form:"i"`
		I37 int `form:"i"`
		I38 int `form:"i"`
		I39 int `form:"i"`
		I40 int `form:"i"`
		I41 int `form:"i"`
		I42 int `form:"i"`
		I43 int `form:"i"`
		I44 int `form:"i"`
		I45 int `form:"i"`
		I46 int `form:"i"`
		I47 int `form:"i"`
		I48 int `form:"i"`
		I49 int `form:"i"`
		I50 int `form:"i"`
		I51 int `form:"i"`
		I52 int `form:"i"`
		I53 int `form:"i"`
		I54 int `form:"i"`
		I55 int `form:"i"`
		I56 int `form:"i"`
		I57 int `form:"i"`
		I58 int `form:"i"`
		I59 int `form:"i"`
		I60 int `form:"i"`
		I61 int `form:"i"`
		I62 int `form:"i"`
		I63 int `form:"i"`
		I64 int `form:"i"`
		I65 int `form:"i"`
		I66 int `form:"i"`
		I67 int `form:"i"`
		I68 int `form:"i"`
		I69 int `form:"i"`
		I70 int `form:"i"`
		I71 int `form:"i"`
		I72 int `form:"i"`
		I73 int `form:"i"`
		I74 int `form:"i"`
		I75 int `form:"i"`
		I76 int `form:"i"`
		I77 int `form:"i"`
		I78 int `form:"i"`
		I79 int `form:"i"`
		I80 int `form:"i"`
		I81 int `form:"i"`
		I82 int `form:"i"`
		I83 int `form:"i"`
		I84 int `form:"i"`
		I85 int `form:"i"`
		I86 int `form:"i"`
		I87 int `form:"i"`
		I88 int `form:"i"`
		I89 int `form:"i"`
		I90 int `form:"i"`
		I91 int `form:"i"`
		I92 int `form:"i"`
		I93 int `form:"i"`
		I94 int `form:"i"`
		I95 int `form:"i"`
		I96 int `form:"i"`
	}

	httpReq, _ := http.NewRequest("GET", "", nil)
	params := make(url.Values)
	params.Add("i", "1")
	params.Add("j", "1")
	params.Add("k", "1")
	for i := 0; i < 100; i++ {
		params.Add(fmt.Sprintf("h%d", i), fmt.Sprint(i))
	}
	httpReq.URL.RawQuery = params.Encode()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		binder := ginBinder.Form
		err := binder.Bind(httpReq, &v)
		if err != nil {
			b.Error("error")
		}
		if v.I != 1 {
			b.FailNow()
		}
	}
}

/***** SmallSlice ******/
type SmallSliceReq struct {
	I []int `query:"i"`
	J []int `query:"j"`
	K []int `query:"k"`
}

//	func getSmallSliceBenchCtx() Ctx {
//		app := New()
//		ctx := app.NewCtx(&fasthttp.RequestCtx{}).(*DefaultCtx)
//		var u = fasthttp.URI{}
//		u.SetQueryString("i=1&i=2&j=1&j=2&k=1&k=2")
//		ctx.Request().SetURI(&u)
//
//		return ctx
//	}
func Benchmark_HertzSmallSlice1(b *testing.B) {
	ctx := app.NewContext(0)
	ctx.URI().SetQueryString("i=1&i=2&j=1&j=2&k=1&k=2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := SmallSliceReq{}
		err := ctx.Bind(&v)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

func Benchmark_GoTagexprSmallSlice1(b *testing.B) {
	ctx := app.NewContext(0)
	ctx.URI().SetQueryString("i=1&i=2&j=1&j=2&k=1&k=2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := SmallSliceReq{}
		err := go_tagexpr.NewBinder().BindAndValidate(&ctx.Request, &v, nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

//func Benchmark_FiberSmallSlice1(b *testing.B) {
//	ctx := getSmallSliceBenchCtx()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		var v = SmallSliceReq{}
//		err := ctx.Bind().Req(&v).Err()
//		if err != nil {
//			b.Error(err)
//			b.FailNow()
//		}
//		if len(v.I) != 2 {
//			b.FailNow()
//		}
//	}
//}

func Benchmark_GinSmallSlice1(b *testing.B) {
	type Req struct {
		I []int `form:"i"`
		J []int `form:"j"`
		K []int `form:"k"`
	}
	httpReq, _ := http.NewRequest("GET", "", nil)
	params := make(url.Values)
	params.Add("i", "1")
	params.Add("i", "1")
	params.Add("j", "1")
	params.Add("j", "1")
	params.Add("k", "1")
	params.Add("k", "1")
	httpReq.URL.RawQuery = params.Encode()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		binder := ginBinder.Form
		err := binder.Bind(httpReq, &v)
		if err != nil {
			b.Error("error")
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

/***** BigQuerySmallSlice ******/
func Benchmark_HertzBigQuerySmallSlice(b *testing.B) {
	ctx := app.NewContext(0)
	str := strings.Builder{}
	for i := 0; i < 100; i++ {
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
	}
	ctx.URI().SetQueryString(str.String() + "i=1&i=2&j=1&j=2&k=1&k=2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := SmallSliceReq{}
		err := ctx.Bind(&v)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

func Benchmark_GoTagexprBigQuerySmallSlice(b *testing.B) {
	ctx := app.NewContext(0)
	str := strings.Builder{}
	for i := 0; i < 100; i++ {
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
	}
	ctx.URI().SetQueryString(str.String() + "i=1&i=2&j=1&j=2&k=1&k=2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := SmallSliceReq{}
		err := go_tagexpr.NewBinder().BindAndValidate(&ctx.Request, &v, nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

//func getBigQuerySmallSliceBenchCtx() Ctx {
//	app := New()
//	ctx := app.NewCtx(&fasthttp.RequestCtx{}).(*DefaultCtx)
//	var u = fasthttp.URI{}
//	str := strings.Builder{}
//	for i := 0; i < 100; i++ {
//		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
//		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
//	}
//	u.SetQueryString(str.String() + "i=1&i=2&j=1&j=2&k=1&k=2")
//	ctx.Request().SetURI(&u)
//
//	return ctx
//}

//func Benchmark_FiberBigQuerySmallSlice(b *testing.B) {
//	ctx := getBigQuerySmallSliceBenchCtx()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		var v = SmallSliceReq{}
//		err := ctx.Bind().Req(&v).Err()
//		if err != nil {
//			b.Error(err)
//			b.FailNow()
//		}
//		if len(v.I) != 2 {
//			b.FailNow()
//		}
//	}
//}

func Benchmark_GinBigQuerySmallSlice(b *testing.B) {
	type Req struct {
		I []int `form:"i"`
		J []int `form:"j"`
		K []int `form:"k"`
	}
	httpReq, _ := http.NewRequest("GET", "", nil)
	params := make(url.Values)
	params.Add("i", "1")
	params.Add("i", "2")
	params.Add("j", "1")
	params.Add("j", "2")
	params.Add("k", "1")
	params.Add("k", "2")
	for i := 0; i < 100; i++ {
		params.Add(fmt.Sprintf("h%d", i), fmt.Sprint(i))
		params.Add(fmt.Sprintf("h%d", i), fmt.Sprint(i))
	}
	httpReq.URL.RawQuery = params.Encode()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := Req{}
		binder := ginBinder.Form
		err := binder.Bind(httpReq, &v)
		if err != nil {
			b.Error("error")
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

/***** BigQueryBigSlice ******/

type BigSliceFieldReq struct {
	I   []int `query:"i"`
	J   []int `query:"j"`
	K   []int `query:"k"`
	I0  []int `query:"i"`
	I1  []int `query:"i"`
	I2  []int `query:"i"`
	I3  []int `query:"i"`
	I4  []int `query:"i"`
	I5  []int `query:"i"`
	I6  []int `query:"i"`
	I7  []int `query:"i"`
	I8  []int `query:"i"`
	I9  []int `query:"i"`
	I10 []int `query:"i"`
	I11 []int `query:"i"`
	I12 []int `query:"i"`
	I13 []int `query:"i"`
	I14 []int `query:"i"`
	I15 []int `query:"i"`
	I16 []int `query:"i"`
	I17 []int `query:"i"`
	I18 []int `query:"i"`
	I19 []int `query:"i"`
	I20 []int `query:"i"`
	I21 []int `query:"i"`
	I22 []int `query:"i"`
	I23 []int `query:"i"`
	I24 []int `query:"i"`
	I25 []int `query:"i"`
	I26 []int `query:"i"`
	I27 []int `query:"i"`
	I28 []int `query:"i"`
	I29 []int `query:"i"`
	I30 []int `query:"i"`
	I31 []int `query:"i"`
	I32 []int `query:"i"`
	I33 []int `query:"i"`
	I34 []int `query:"i"`
	I35 []int `query:"i"`
	I36 []int `query:"i"`
	I37 []int `query:"i"`
	I38 []int `query:"i"`
	I39 []int `query:"i"`
	I40 []int `query:"i"`
	I41 []int `query:"i"`
	I42 []int `query:"i"`
	I43 []int `query:"i"`
	I44 []int `query:"i"`
	I45 []int `query:"i"`
	I46 []int `query:"i"`
	I47 []int `query:"i"`
	I48 []int `query:"i"`
	I49 []int `query:"i"`
	I50 []int `query:"i"`
	I51 []int `query:"i"`
	I52 []int `query:"i"`
	I53 []int `query:"i"`
	I54 []int `query:"i"`
	I55 []int `query:"i"`
	I56 []int `query:"i"`
	I57 []int `query:"i"`
	I58 []int `query:"i"`
	I59 []int `query:"i"`
	I60 []int `query:"i"`
	I61 []int `query:"i"`
	I62 []int `query:"i"`
	I63 []int `query:"i"`
	I64 []int `query:"i"`
	I65 []int `query:"i"`
	I66 []int `query:"i"`
	I67 []int `query:"i"`
	I68 []int `query:"i"`
	I69 []int `query:"i"`
	I70 []int `query:"i"`
	I71 []int `query:"i"`
	I72 []int `query:"i"`
	I73 []int `query:"i"`
	I74 []int `query:"i"`
	I75 []int `query:"i"`
	I76 []int `query:"i"`
	I77 []int `query:"i"`
	I78 []int `query:"i"`
	I79 []int `query:"i"`
	I80 []int `query:"i"`
	I81 []int `query:"i"`
	I82 []int `query:"i"`
	I83 []int `query:"i"`
	I84 []int `query:"i"`
	I85 []int `query:"i"`
	I86 []int `query:"i"`
	I87 []int `query:"i"`
	I88 []int `query:"i"`
	I89 []int `query:"i"`
	I90 []int `query:"i"`
	I91 []int `query:"i"`
	I92 []int `query:"i"`
	I93 []int `query:"i"`
	I94 []int `query:"i"`
	I95 []int `query:"i"`
	I96 []int `query:"i"`
}

func Benchmark_HertzBigQueryBigSlice(b *testing.B) {
	ctx := app.NewContext(0)
	str := strings.Builder{}
	for i := 0; i < 100; i++ {
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
	}
	ctx.URI().SetQueryString(str.String() + "i=1&i=2&j=1&j=2&k=1&k=2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := BigSliceFieldReq{}
		err := ctx.Bind(&v)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

func Benchmark_GoTagexprBigQueryBigSlice(b *testing.B) {
	ctx := app.NewContext(0)
	str := strings.Builder{}
	for i := 0; i < 100; i++ {
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
		str.WriteString(fmt.Sprintf("h%d=%d&", i, i))
	}
	ctx.URI().SetQueryString(str.String() + "i=1&i=2&j=1&j=2&k=1&k=2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := BigSliceFieldReq{}
		err := go_tagexpr.NewBinder().BindAndValidate(&ctx.Request, &v, nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

//func Benchmark_FiberBigQueryBigSlice(b *testing.B) {
//	ctx := getBigQuerySmallSliceBenchCtx()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		var v = BigSliceFieldReq{}
//		err := ctx.Bind().Req(&v).Err()
//		if err != nil {
//			b.Error(err)
//			b.FailNow()
//		}
//		if len(v.I) != 2 {
//			b.FailNow()
//		}
//	}
//}

func Benchmark_GinBigQueryBigSlice(b *testing.B) {
	type BigSliceFieldReq struct {
		I   []int `form:"i"`
		J   []int `form:"j"`
		K   []int `form:"k"`
		I0  []int `form:"i"`
		I1  []int `form:"i"`
		I2  []int `form:"i"`
		I3  []int `form:"i"`
		I4  []int `form:"i"`
		I5  []int `form:"i"`
		I6  []int `form:"i"`
		I7  []int `form:"i"`
		I8  []int `form:"i"`
		I9  []int `form:"i"`
		I10 []int `form:"i"`
		I11 []int `form:"i"`
		I12 []int `form:"i"`
		I13 []int `form:"i"`
		I14 []int `form:"i"`
		I15 []int `form:"i"`
		I16 []int `form:"i"`
		I17 []int `form:"i"`
		I18 []int `form:"i"`
		I19 []int `form:"i"`
		I20 []int `form:"i"`
		I21 []int `form:"i"`
		I22 []int `form:"i"`
		I23 []int `form:"i"`
		I24 []int `form:"i"`
		I25 []int `form:"i"`
		I26 []int `form:"i"`
		I27 []int `form:"i"`
		I28 []int `form:"i"`
		I29 []int `form:"i"`
		I30 []int `form:"i"`
		I31 []int `form:"i"`
		I32 []int `form:"i"`
		I33 []int `form:"i"`
		I34 []int `form:"i"`
		I35 []int `form:"i"`
		I36 []int `form:"i"`
		I37 []int `form:"i"`
		I38 []int `form:"i"`
		I39 []int `form:"i"`
		I40 []int `form:"i"`
		I41 []int `form:"i"`
		I42 []int `form:"i"`
		I43 []int `form:"i"`
		I44 []int `form:"i"`
		I45 []int `form:"i"`
		I46 []int `form:"i"`
		I47 []int `form:"i"`
		I48 []int `form:"i"`
		I49 []int `form:"i"`
		I50 []int `form:"i"`
		I51 []int `form:"i"`
		I52 []int `form:"i"`
		I53 []int `form:"i"`
		I54 []int `form:"i"`
		I55 []int `form:"i"`
		I56 []int `form:"i"`
		I57 []int `form:"i"`
		I58 []int `form:"i"`
		I59 []int `form:"i"`
		I60 []int `form:"i"`
		I61 []int `form:"i"`
		I62 []int `form:"i"`
		I63 []int `form:"i"`
		I64 []int `form:"i"`
		I65 []int `form:"i"`
		I66 []int `form:"i"`
		I67 []int `form:"i"`
		I68 []int `form:"i"`
		I69 []int `form:"i"`
		I70 []int `form:"i"`
		I71 []int `form:"i"`
		I72 []int `form:"i"`
		I73 []int `form:"i"`
		I74 []int `form:"i"`
		I75 []int `form:"i"`
		I76 []int `form:"i"`
		I77 []int `form:"i"`
		I78 []int `form:"i"`
		I79 []int `form:"i"`
		I80 []int `form:"i"`
		I81 []int `form:"i"`
		I82 []int `form:"i"`
		I83 []int `form:"i"`
		I84 []int `form:"i"`
		I85 []int `form:"i"`
		I86 []int `form:"i"`
		I87 []int `form:"i"`
		I88 []int `form:"i"`
		I89 []int `form:"i"`
		I90 []int `form:"i"`
		I91 []int `form:"i"`
		I92 []int `form:"i"`
		I93 []int `form:"i"`
		I94 []int `form:"i"`
		I95 []int `form:"i"`
		I96 []int `form:"i"`
	}

	httpReq, _ := http.NewRequest("GET", "", nil)
	params := make(url.Values)
	params.Add("i", "1")
	params.Add("i", "2")
	params.Add("j", "1")
	params.Add("j", "2")
	params.Add("k", "1")
	params.Add("k", "2")
	for i := 0; i < 100; i++ {
		params.Add(fmt.Sprintf("h%d", i), fmt.Sprint(i))
		params.Add(fmt.Sprintf("h%d", i), fmt.Sprint(i))
	}
	httpReq.URL.RawQuery = params.Encode()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := BigSliceFieldReq{}
		binder := ginBinder.Form
		err := binder.Bind(httpReq, &v)
		if err != nil {
			b.Error("error")
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

/***** SmallQueryBigSlice ******/

func Benchmark_HertzSmallQueryBigSlice(b *testing.B) {
	ctx := app.NewContext(0)
	ctx.URI().SetQueryString("i=1&i=2&j=1&j=2&k=1&k=2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := BigSliceFieldReq{}
		err := ctx.Bind(&v)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

func Benchmark_GoTagexprSmallQueryBigSlice(b *testing.B) {
	ctx := app.NewContext(0)
	ctx.URI().SetQueryString("i=1&i=2&j=1&j=2&k=1&k=2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := BigSliceFieldReq{}
		err := go_tagexpr.NewBinder().BindAndValidate(&ctx.Request, &v, nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}

//func Benchmark_FiberSmallQueryBigSlice(b *testing.B) {
//	ctx := getSmallSliceBenchCtx()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		var v = BigSliceFieldReq{}
//		err := ctx.Bind().Req(&v).Err()
//		if err != nil {
//			b.Error(err)
//			b.FailNow()
//		}
//		if len(v.I) != 2 {
//			b.FailNow()
//		}
//	}
//}

func Benchmark_GinSmallQueryBigSlice(b *testing.B) {
	type BigSliceFieldReq struct {
		I   []int `form:"i"`
		J   []int `form:"j"`
		K   []int `form:"k"`
		I0  []int `form:"i"`
		I1  []int `form:"i"`
		I2  []int `form:"i"`
		I3  []int `form:"i"`
		I4  []int `form:"i"`
		I5  []int `form:"i"`
		I6  []int `form:"i"`
		I7  []int `form:"i"`
		I8  []int `form:"i"`
		I9  []int `form:"i"`
		I10 []int `form:"i"`
		I11 []int `form:"i"`
		I12 []int `form:"i"`
		I13 []int `form:"i"`
		I14 []int `form:"i"`
		I15 []int `form:"i"`
		I16 []int `form:"i"`
		I17 []int `form:"i"`
		I18 []int `form:"i"`
		I19 []int `form:"i"`
		I20 []int `form:"i"`
		I21 []int `form:"i"`
		I22 []int `form:"i"`
		I23 []int `form:"i"`
		I24 []int `form:"i"`
		I25 []int `form:"i"`
		I26 []int `form:"i"`
		I27 []int `form:"i"`
		I28 []int `form:"i"`
		I29 []int `form:"i"`
		I30 []int `form:"i"`
		I31 []int `form:"i"`
		I32 []int `form:"i"`
		I33 []int `form:"i"`
		I34 []int `form:"i"`
		I35 []int `form:"i"`
		I36 []int `form:"i"`
		I37 []int `form:"i"`
		I38 []int `form:"i"`
		I39 []int `form:"i"`
		I40 []int `form:"i"`
		I41 []int `form:"i"`
		I42 []int `form:"i"`
		I43 []int `form:"i"`
		I44 []int `form:"i"`
		I45 []int `form:"i"`
		I46 []int `form:"i"`
		I47 []int `form:"i"`
		I48 []int `form:"i"`
		I49 []int `form:"i"`
		I50 []int `form:"i"`
		I51 []int `form:"i"`
		I52 []int `form:"i"`
		I53 []int `form:"i"`
		I54 []int `form:"i"`
		I55 []int `form:"i"`
		I56 []int `form:"i"`
		I57 []int `form:"i"`
		I58 []int `form:"i"`
		I59 []int `form:"i"`
		I60 []int `form:"i"`
		I61 []int `form:"i"`
		I62 []int `form:"i"`
		I63 []int `form:"i"`
		I64 []int `form:"i"`
		I65 []int `form:"i"`
		I66 []int `form:"i"`
		I67 []int `form:"i"`
		I68 []int `form:"i"`
		I69 []int `form:"i"`
		I70 []int `form:"i"`
		I71 []int `form:"i"`
		I72 []int `form:"i"`
		I73 []int `form:"i"`
		I74 []int `form:"i"`
		I75 []int `form:"i"`
		I76 []int `form:"i"`
		I77 []int `form:"i"`
		I78 []int `form:"i"`
		I79 []int `form:"i"`
		I80 []int `form:"i"`
		I81 []int `form:"i"`
		I82 []int `form:"i"`
		I83 []int `form:"i"`
		I84 []int `form:"i"`
		I85 []int `form:"i"`
		I86 []int `form:"i"`
		I87 []int `form:"i"`
		I88 []int `form:"i"`
		I89 []int `form:"i"`
		I90 []int `form:"i"`
		I91 []int `form:"i"`
		I92 []int `form:"i"`
		I93 []int `form:"i"`
		I94 []int `form:"i"`
		I95 []int `form:"i"`
		I96 []int `form:"i"`
	}

	httpReq, _ := http.NewRequest("GET", "", nil)
	params := make(url.Values)
	params.Add("i", "1")
	params.Add("i", "2")
	params.Add("j", "1")
	params.Add("j", "2")
	params.Add("k", "1")
	params.Add("k", "2")
	httpReq.URL.RawQuery = params.Encode()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := BigSliceFieldReq{}
		binder := ginBinder.Form
		err := binder.Bind(httpReq, &v)
		if err != nil {
			b.Error("error")
		}
		if len(v.I) != 2 {
			b.FailNow()
		}
	}
}
