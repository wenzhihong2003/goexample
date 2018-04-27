package main

import (
	"fmt"
	"net"

	"github.com/wenzhihong2003/goexample/go-grpc-example/proto/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	Address = "127.0.0.1:50052"
)

type helloService struct{}

func (helloService) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "无Token 认证信息")
	}
	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return nil, status.Errorf(codes.Unauthenticated, "Token认证信息无效")
	}
	resp := new(hello.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s.\n Token info: appid=%s, appkey=%s", in.Name, appid, appkey)
	return resp, nil
}

var HelloService = helloService{}

func main() {
	listener, e := net.Listen("tcp", Address)
	if e != nil {
		grpclog.Fatalf("failed to listen: %v", e)
	}

	creds, e := credentials.NewServerTLSFromFile("../../keys/mytestserver.pem", "../../keys/mytestserver.key")
	if e != nil {
		grpclog.Fatalf("Failed to generate credentials %v", e)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	hello.RegisterHelloServer(s, HelloService)
	grpclog.Infoln("Listen on " + Address + " with TLS + TOKEN")
	s.Serve(listener)
}
