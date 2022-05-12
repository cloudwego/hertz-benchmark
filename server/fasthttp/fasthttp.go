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

package main

import (
	"unsafe"

	"github.com/cloudwego/hertz-benchmark/perf"
	"github.com/cloudwego/hertz-benchmark/runner"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

const (
	port        = ":8002"
	debugPort   = ":18002"
	actionQuery = "action"
)

var recorder = perf.NewRecorder("FastHttp@Server")

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(debugPort)
	}()

	r := router.New()

	r.POST("/", requestHandler)

	s := &fasthttp.Server{
		Handler: r.Handler,
	}

	s.ListenAndServe(port)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	runner.ProcessRequest(recorder, b2s(ctx.QueryArgs().Peek(actionQuery)))

	ctx.SetContentType("text/plain; charset=utf8")
	ctx.Response.SetBody(ctx.Request.Body())
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
