package main

import (
	"os"

	"github.com/golang/glog"
	"github.com/wenzhihong2003/goexample/go-grpc-example/proto/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	_ "google.golang.org/grpc/grpclog/glogger"
)

const (
	Address = "127.0.0.1:50052"
)

func init() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
}

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	client := hello.NewHelloClient(conn)

	// 调用方法
	reqBody := new(hello.HelloRequest)
	reqBody.Name = "gRPC"
	response, err := client.SayHello(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}
	grpclog.Infoln(response.Message)

	glog.Flush() // 要自己调用 flush 保证把内容都刷到文件里去. 不然的话, 可能文件里会没有内容. glog是定时刷的, 默认为10s
}
