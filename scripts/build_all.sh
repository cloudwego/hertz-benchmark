#!/bin/bash

GOEXEC=${GOEXEC:-"go"}

# clean
rm -rf output/ && mkdir -p output/bin/ && mkdir -p output/log/

$GOEXEC mod tidy && go mod verify

# build clients
$GOEXEC build -v -o output/bin/hertz_client ./client/hertz/hertz.go
$GOEXEC build -v -o output/bin/net_client ./client/net/net.go
$GOEXEC build -v -o output/bin/fasthttp_client ./client/fasthttp/fasthttp.go

# build servers
$GOEXEC build -v -o output/bin/hertz_server ./server/hertz/hertz.go
$GOEXEC build -v -o output/bin/gin_server ./server/gin/gin.go
$GOEXEC build -v -o output/bin/fasthttp_server ./server/fasthttp/fasthttp.go
$GOEXEC build -v -o output/bin/fiber_server ./server/fiber/fiber.go
$GOEXEC build -v -o output/bin/hex_server ./server/hex/.
