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
	"fmt"

	"github.com/cloudwego/hertz-benchmark/runner"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/netpoll"
)

func main() {
	runner.Main("hertz_client", NewHertzClient)
}

type Client struct {
	c *client.Client
}

func NewHertzClient(opt *runner.Options) runner.Client {
	_ = netpoll.SetNumLoops(2)
	client, err := client.NewClient(client.WithMaxConnsPerHost(1000))
	if err != nil {
		panic(err)
	}
	return &Client{
		c: client,
	}
}

func (c *Client) Echo(action, uri, body, header string) error {
	req, resp := protocol.AcquireRequest(), protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseResponse(resp)
		protocol.ReleaseRequest(req)
	}()
	req.SetMethod(consts.MethodPost)
	req.Header.Set("test-header", header)
	req.SetBodyString(body)
	req.SetRequestURI(fmt.Sprintf("%s?action=%s", uri, action))
	ctx := context.Background()
	err := c.c.Do(ctx, req, resp)
	return err
}
