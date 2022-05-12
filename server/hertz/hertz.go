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
	"context"
	"time"
	"unsafe"

	"github.com/cloudwego/hertz-benchmark/perf"
	"github.com/cloudwego/hertz-benchmark/runner"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/netpoll"
)

const (
	port        = ":8001"
	debugPort   = ":18001"
	actionQuery = "action"
)

var recorder = perf.NewRecorder("Hertz@Server")

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(debugPort)
	}()

	netpoll.SetNumLoops(2)
	opts := []config.Option{
		server.WithHostPorts(port),
		server.WithIdleTimeout(time.Second * 10),
		server.WithReadTimeout(time.Second * 10),
	}
	h := server.New(opts...)

	h.POST("/", echoHandler)

	h.Spin()
}

func echoHandler(c context.Context, ctx *app.RequestContext) {
	runner.ProcessRequest(recorder, b2s(ctx.QueryArgs().Peek(actionQuery)))
	ctx.SetContentType("text/plain; charset=utf8")
	ctx.Response.SetBody(ctx.Request.Body())
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
