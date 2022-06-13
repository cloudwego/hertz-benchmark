/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package runner

import (
	"flag"
	"fmt"
	"log"

	"github.com/cloudwego/hertz-benchmark/perf"
)

var (
	address    string
	bodySize   int
	headerSize int
	total      int64
	concurrent int
	poolSize   int
	sleepTime  int
	server     string
)

type Options struct {
	Address  string
	Body     []byte
	PoolSize int
}

type ClientNewer func(opt *Options) Client

type Client interface {
	Echo(action, uri, body, header string) (err error)
}

func initFlags() {
	flag.StringVar(&address, "addr", "", "client call address")
	flag.IntVar(&bodySize, "b", 1024, "body size once")
	flag.IntVar(&headerSize, "h", 0, "header size once")
	flag.IntVar(&concurrent, "c", 1, "call concurrent")
	flag.Int64Var(&total, "n", 1, "call total nums")
	flag.IntVar(&poolSize, "pool", 10, "conn poll size")
	flag.IntVar(&sleepTime, "sleep", 0, "sleep time for every request handler")
	flag.StringVar(&server, "s", "", "call server")
	flag.Parse()
}

func Main(name string, newer ClientNewer) {
	initFlags()

	r := NewRunner()

	opt := &Options{
		Address:  address,
		PoolSize: poolSize,
	}
	cli := newer(opt)
	body := string(make([]byte, bodySize))
	headerSlice := make([]byte, headerSize)
	for i := 0; i < len(headerSlice); i++ {
		headerSlice[i] = 'a'
	}
	header := string(headerSlice)
	action := EchoAction
	handler := func() error { return cli.Echo(action, address, body, header) }
	// === warming ===
	r.Warmup(handler, concurrent, 100*1000)
	// === beginning ===
	if err := cli.Echo(BeginAction, address, "", ""); err != nil {
		log.Fatalf("beginning server failed: %v", err)
	}
	recorder := perf.NewRecorder(fmt.Sprintf("%s@Client", name))
	recorder.Begin()

	// === benching ===
	r.Run(fmt.Sprintf("%s - %s", name, server), handler, bodySize, concurrent, total, bodySize, headerSize)

	// == ending ===
	recorder.End()
	if err := cli.Echo(EndAction, address, "", ""); err != nil {
		log.Fatalf("ending server failed: %v", err)
	}
	// === reporting ===
	recorder.Report() // report client
	fmt.Printf("\n\n")
}
