package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/gopherguides/training/distributed-systems/grpc/src/hoover/hoover"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
	url     = "http://www.google.com"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer the close
	defer conn.Close()

	// create our service with the connection
	c := hoover.NewServiceClient(conn)

	r, err := c.Get(context.Background(), &hoover.GetRequest{Url: url})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// print the body of the page retrieved
	fmt.Println(r.Body)

	// print response code and elapsed time
	log.Printf("Response Code: %d, Elapsed: %s", r.ResponseCode, convertDuration(r.Elapsed))
}

// convert a protubuff duration to a go duration
func convertDuration(d *duration.Duration) time.Duration {
	return time.Duration((d.Seconds * int64(time.Second)) + int64(d.Nanos))
}
