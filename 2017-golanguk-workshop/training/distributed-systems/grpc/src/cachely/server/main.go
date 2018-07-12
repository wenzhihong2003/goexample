package main

import (
	"log"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/gopherguides/training/distributed-systems/grpc/src/cachely/cachely"

	"golang.org/x/sync/syncmap"
)

type server struct {
	data syncmap.Map
}

func (s *server) Get(_ context.Context, req *cachely.GetRequest) (*cachely.GetResponse, error) {
	key := req.GetKey()
	log.Printf("looking up key %q\n", key)
	if v, ok := s.data.Load(key); ok {
		log.Printf("found key %q\n", key)
		return &cachely.GetResponse{
			Key:   key,
			Value: v.([]byte),
		}, nil
	}
	log.Printf("key not found %q\n", key)
	return nil, grpc.Errorf(codes.NotFound, "could not find key %s", key)
}

func (s *server) Put(_ context.Context, req *cachely.PutRequest) (*cachely.PutResponse, error) {
	log.Printf("storing key %q\n", req.GetKey())
	s.data.Store(req.GetKey(), req.GetValue())
	return &cachely.PutResponse{
		Key: req.GetKey(),
	}, nil
}

func (s *server) Delete(_ context.Context, req *cachely.DeleteRequest) (*cachely.DeleteResponse, error) {
	log.Printf("deleting key %q\n", req.GetKey())
	s.data.Delete(req.GetKey())
	return &cachely.DeleteResponse{
		Key: req.GetKey(),
	}, nil
}

func main() {
	// open a port to communicate on
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new grpc server
	s := grpc.NewServer()

	// register our service
	cachely.RegisterCacheServer(s, &server{
		data: syncmap.Map{},
	})

	// start listening and responding
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
