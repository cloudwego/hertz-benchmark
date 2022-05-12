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
	"bytes"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"
	"unsafe"

	"github.com/cloudwego/hertz-benchmark/runner"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	runner.Main("net_client", NewNetClient)
}

type Client struct {
	c *http.Client
}

func NewNetClient(opt *runner.Options) runner.Client {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
		},
	}
	return &Client{
		c: client,
	}
}

func (c *Client) Echo(action, uri, body, header string) error {
	reqBody := bytes.NewBuffer(s2b(body))
	req, err := http.NewRequest(consts.MethodPost, fmt.Sprintf("%s?action=%s", uri, action), reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("test-header", header)
	resp, err := c.c.Do(req)
	if resp.Body != nil {
		resp.Body.Close()
	}
	return err
}

func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
