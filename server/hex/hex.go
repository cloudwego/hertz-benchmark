package main

import (
	"github.com/cloudwego/hertz-benchmark/perf"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net"

	"github.com/cloudwego/hertz-benchmark/server/hex/kitex_gen/hello/example/helloservice"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/cloudwego/kitex/server"
)

func main() {

	// start pprof server
	go func() {
		err := perf.ServeMonitor(":18005")
		if err != nil {
			hlog.Error(err)
		}
	}()

	opts := kitexInit()

	svr := helloservice.NewServer(new(HelloServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	opts = append(opts, server.
		WithTransHandlerFactory(&mixTransHandlerFactory{nil}))

	// address
	addr, err := net.ResolveTCPAddr("tcp", ":8005")
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	return
}
