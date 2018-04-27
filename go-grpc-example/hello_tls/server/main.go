package main

import (
	"net"
	"os"

	"github.com/wenzhihong2003/goexample/go-grpc-example/proto/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	Address = "127.0.0.1:50052"
)

type helloService struct {
}

func (h helloService) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	resp := new(hello.HelloResponse)
	resp.Message = "Hello " + in.Name + "!!!"
	return resp, nil
}

var HelloService = helloService{}

func init() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
}

func main() {
	listener, er := net.Listen("tcp", Address)
	if er != nil {
		grpclog.Fatalf("failed to listen: %v", er)
	}

	// TLS 认证
	creds, er := credentials.NewServerTLSFromFile("../../keys/mytestserver.pem","../../keys/mytestserver.key")
	if er != nil {
		grpclog.Fatalf("Failed to generate credentials %v", er)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	hello.RegisterHelloServer(s, HelloService)
	grpclog.Infoln("Listen on " + Address + " with TLS")

	s.Serve(listener)
}
