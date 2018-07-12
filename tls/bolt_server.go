package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
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

	config1 := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// 因为需要验证客户端，我们需要额外配置下面两个字段
		// 然后客户端也配置这个clientCertPool
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  clientCertPool,
	}
	ln, err := tls.Listen("tcp", ":443", config1)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}

	}
}
