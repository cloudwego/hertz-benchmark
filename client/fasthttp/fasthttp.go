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
	"fmt"

	"github.com/cloudwego/hertz-benchmark/runner"
	"github.com/valyala/fasthttp"
)

func main() {
	runner.Main("fasthttp_client", NewFasthttpClient)
}

type Client struct {
	c *fasthttp.Client
}

func NewFasthttpClient(opt *runner.Options) runner.Client {
	client := &fasthttp.Client{
		MaxConnsPerHost: 1000,
	}
	return &Client{
		c: client,
	}
}

func (c *Client) Echo(action, uri, body, header string) error {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	req.Header.SetMethod("POST")
	req.Header.Set("test-header", header)
	req.SetBodyString(body)
	req.SetRequestURI(fmt.Sprintf("%s?action=%s", uri, action))
	err := c.c.Do(req, resp)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
