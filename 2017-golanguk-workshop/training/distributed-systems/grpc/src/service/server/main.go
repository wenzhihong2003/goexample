package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/gopherguides/training/distributed-systems/grpc/src/service/discovery"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type services struct {
	serviceType discovery.ServiceType
	nodes       map[string]discovery.Node
}

func NewServices() *services {
	return &services{
		nodes: make(map[string]discovery.Node),
	}
}

func (s *services) toService(name string) *discovery.Service {
	ds := &discovery.Service{
		Name: name,
		Type: s.serviceType,
	}
	for _, n := range s.nodes {
		ds.Nodes = append(ds.Nodes, &n)
	}

	return ds
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	services map[string]*services
}

func NewServer() *server {
	return &server{
		services: make(map[string]*services),
	}
}

// Register will register the service
func (s *server) Register(ctx context.Context, in *discovery.RegistrationRequest) (*discovery.RegistrationReply, error) {
	if s.services[in.Name] == nil {
		s.services[in.Name] = NewServices()
		s.services[in.Name].serviceType = in.Type
	}
	// make sure they don't mix up service types
	if in.Type != s.services[in.Name].serviceType {
		return nil, fmt.Errorf("mismatched service types. have %d, given %d", s.services[in.Name].serviceType, in.Type)
	}
	now := time.Now()

	s.services[in.Name].nodes[in.UUID] = discovery.Node{
		UUID:      in.UUID,
		URI:       in.URI,
		Leader:    in.Leader,
		UpdatedAt: timeToTimestamp(now),
		Heartbeat: timeToTimestamp(now),
	}

	return &discovery.RegistrationReply{Success: true}, nil
}

// Heartbeat will update the service for a last heartbeat metric
func (s *server) Heartbeat(ctx context.Context, in *discovery.HeartbeatRequest) (*discovery.HeartbeatReply, error) {
	return &discovery.HeartbeatReply{Success: true}, nil
}

func (s *server) List(ctx context.Context, in *discovery.ListRequest) (*discovery.ListReply, error) {
	if in.Name != "" {
		if s.services[in.Name] == nil {
			return &discovery.ListReply{}, nil
		}
		out := &discovery.ListReply{
			Services: make(map[string]*discovery.Service),
		}
		out.Services[in.Name] = s.services[in.Name].toService(in.Name)
		return out, nil
	}

	out := &discovery.ListReply{
		Services: make(map[string]*discovery.Service),
	}
	for key := range s.services {
		out.Services[key] = s.services[key].toService(key)
	}

	return out, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	discovery.RegisterDiscoveryServer(s, NewServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func timeToTimestamp(t time.Time) *timestamp.Timestamp {
	seconds := t.Unix()
	nanos := t.UnixNano() - (seconds * int64(time.Second))
	return &timestamp.Timestamp{
		Seconds: seconds,
		Nanos:   int32(nanos),
	}
}
