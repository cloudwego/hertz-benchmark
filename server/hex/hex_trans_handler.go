package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz-benchmark/perf"
	"github.com/cloudwego/hertz-benchmark/runner"
	"github.com/cloudwego/kitex/pkg/remote/trans/gonet"
	"net"
	"regexp"
	"unsafe"

	"github.com/cloudwego/hertz/pkg/app"
	hertzServer "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/detection"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2"
)

type mixTransHandlerFactory struct {
	originFactory remote.ServerTransHandlerFactory
}

type transHandler struct {
	remote.ServerTransHandler
}

// SetInvokeHandleFunc is used to set invoke handle func.
func (t *transHandler) SetInvokeHandleFunc(inkHdlFunc endpoint.Endpoint) {
	t.ServerTransHandler.(remote.InvokeHandleFuncSetter).SetInvokeHandleFunc(inkHdlFunc)
}

func (m mixTransHandlerFactory) NewTransHandler(opt *remote.ServerOption) (remote.ServerTransHandler, error) {
	var kitexOrigin remote.ServerTransHandler
	var err error

	if m.originFactory != nil {
		kitexOrigin, err = m.originFactory.NewTransHandler(opt)
	} else {
		// if no customized factory just use the default factory under detection pkg.
		kitexOrigin, err = detection.NewSvrTransHandlerFactory(gonet.NewSvrTransHandlerFactory(), nphttp2.NewSvrTransHandlerFactory()).NewTransHandler(opt)
	}
	if err != nil {
		return nil, err
	}
	return &transHandler{ServerTransHandler: kitexOrigin}, nil
}

var httpReg = regexp.MustCompile(`^(?:GET |POST|PUT|DELE|HEAD|OPTI|CONN|TRAC|PATC)$`)

func (t *transHandler) OnRead(ctx context.Context, conn net.Conn) error {
	c, ok := conn.(network.Conn)
	if ok {
		pre, _ := c.Peek(4)
		if httpReg.Match(pre) {
			//klog.Info("using Hertz to process request")
			err := hertzEngine.Serve(ctx, c)
			if err != nil {
				err = errors.New(fmt.Sprintf("HERTZ: %s", err.Error()))
			}
			return err
		}
	}
	return t.ServerTransHandler.OnRead(ctx, conn)
}

func initHertz() *route.Engine {
	h := hertzServer.New(hertzServer.WithIdleTimeout(0))

	h.POST("/", echoHandler)
	err := h.Engine.Init()
	if err != nil {
		panic(err)
	}

	err = h.MarkAsRunning()
	if err != nil {
		panic(err)
	}

	return h.Engine
}

var (
	recorder    = perf.NewRecorder("Hex@Server")
	actionQuery = "action"
)

func echoHandler(c context.Context, ctx *app.RequestContext) {
	runner.ProcessRequest(recorder, b2s(ctx.QueryArgs().Peek(actionQuery)))
	ctx.SetContentType("text/plain; charset=utf8")
	ctx.Response.SetBody(ctx.Request.Body())
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

var hertzEngine *route.Engine

func init() {
	hertzEngine = initHertz()
}
