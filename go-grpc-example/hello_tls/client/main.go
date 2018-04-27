package main

import (
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

func init() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
}

func main() {
	creds, e := credentials.NewClientTLSFromFile("../../keys/mytestserver.pem", "wen")
	if e != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", e)
	}
	conn, e := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	if e != nil {
		grpclog.Fatalln(e)
	}
	defer conn.Close()

	c := hello.NewHelloClient(conn)
	reqBody := new(hello.HelloRequest)
	reqBody.Name = "grpc"
	resp, e := c.SayHello(context.Background(), reqBody)
	if e != nil {
		grpclog.Fatalln(e)
	}
	grpclog.Infoln(resp.Message)

}
