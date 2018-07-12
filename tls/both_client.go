package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		log.Println(err)
		return
	}

	certBytes, err := ioutil.ReadFile("client.pem")
	if err != nil {
		panic("Unable to read cert.pem")
	}

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}

	// InsecureSkipVerify用来控制客户端是否证书和服务器主机名。如果设置为true,则不会校验证书以及证书中的主机名和服务器主机名是否一致。
	// 因为在我们的例子中使用自签名的证书，所以设置它为true,仅仅用于测试目的
	conf := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	fmt.Println(string(buf[:n]))
}
