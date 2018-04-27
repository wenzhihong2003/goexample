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
	// 是否开启TLS
	OpenTLS = true
)

type customCredential struct {}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid": "101010",
		"appkey": "i am key",
	}, nil
}

func (customCredential) RequireTransportSecurity() bool {
	if OpenTLS {
		return true
	}

	return false
}

func init() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
}

func main() {
	var err error
	var opts []grpc.DialOption

	if OpenTLS {
		creds, err := credentials.NewClientTLSFromFile("../../keys/mytestserver.pem", "wen")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}else {
		opts = append(opts, grpc.WithInsecure())
	}

	// 使用自定义认证
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		grpclog.Fatalln(err)
	}

	defer conn.Close()

	c := hello.NewHelloClient(conn)
	reqBody := new(hello.HelloRequest)
	reqBody.Name = "grpc"
	r, err := c.SayHello(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}

	grpclog.Infoln(r.Message)
}
