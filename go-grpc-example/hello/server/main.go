package main

import (
	"net"
	"os"

	"github.com/wenzhihong2003/goexample/go-grpc-example/proto/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address grpc 服务地址
	Address = "127.0.0.1:50052"
)

type helloService struct {
}

func (h helloService) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	resp := new(hello.HelloResponse)
	resp.Message = "Hello " + in.Name + "."
	return resp, nil
}

var HelloService = helloService{}

func init() {
	// 使用 grpclog 包的日志输出. 这个是直接输出到控制台了.
	// 也可考虑使用  google.golang.org/grpc/grpclog/glogger 包, 他是使用 github.com/golang/glog 这个包来进行日志输出的.
	loggerV2 := grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
	grpclog.SetLoggerV2(loggerV2)
}

func main() {
	listener, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	hello.RegisterHelloServer(server, HelloService)
	grpclog.Infoln("Listen on " + Address)

	server.Serve(listener)
}
