package main

import (
	"net"
	"net/http"
	"os"

	"github.com/wenzhihong2003/goexample/go-grpc-example/proto/hello"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
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
	listener, e := net.Listen("tcp", Address)
	if e != nil {
		grpclog.Fatalf("failed to listen: %v", e)
	}

	s:=grpc.NewServer()
	hello.RegisterHelloServer(s, HelloService)

	// 开启 trace
	go startTrace()

	grpclog.Infoln("listen on " + Address)
	s.Serve(listener)
}

func startTrace() {
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
	go http.ListenAndServe(":50051", nil)
	grpclog.Infoln("Trace listen on 50051")
}


