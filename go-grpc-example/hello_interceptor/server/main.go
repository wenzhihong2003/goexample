package main

import (
	"net"
	"os"

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

type helloService struct {
}

func (h helloService) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	resp := new(hello.HelloResponse)
	resp.Message = "Hello " + in.Name + "."
	return resp, nil
}

var HelloSerice = helloService{}

// auth 验证
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "无token认证信息")
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
		return status.Errorf(codes.Unauthenticated, "token认证信息无效")
	}
	return nil
}

func init() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
}

func main() {
	listener, e := net.Listen("tcp", Address)
	if e != nil {
		grpclog.Fatalf("Failed to listen: %v", e)
	}
	var opts []grpc.ServerOption
	creds, e := credentials.NewServerTLSFromFile("../../keys/mytestserver.pem", "../../keys/mytestserver.key")
	if e != nil {
		grpclog.Fatalf("Failed to generate credentials %v", e)
	}
	opts = append(opts, grpc.Creds(creds))
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}
		// 继续处理
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	s := grpc.NewServer(opts...)
	hello.RegisterHelloServer(s, HelloSerice)
	grpclog.Infoln("Listen on " + Address + " with TLS + token + intercept")
	s.Serve(listener)
}
