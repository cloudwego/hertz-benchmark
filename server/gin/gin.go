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
	"io/ioutil"

	"github.com/cloudwego/hertz-benchmark/perf"
	"github.com/cloudwego/hertz-benchmark/runner"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gin-gonic/gin"
)

const (
	port        = ":8003"
	debugPort   = ":18003"
	actionQuery = "action"
)

var recorder = perf.NewRecorder("Gin@Server")

func main() {
	// start pprof server
	go func() {
		err := perf.ServeMonitor(debugPort)
		if err != nil {
			hlog.Error(err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.POST("/", echoHandler)

	err := r.Run(port)
	if err != nil {
		hlog.Error(err)
	}
}

func echoHandler(c *gin.Context) {
	runner.ProcessRequest(recorder, c.Query(actionQuery))
	b, _ := ioutil.ReadAll(c.Request.Body)
	_, err := c.Writer.Write(b)
	if err != nil {
		hlog.Error(err)
	}
}
