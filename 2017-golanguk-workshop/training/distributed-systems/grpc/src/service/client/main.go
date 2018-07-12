package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gopherguides/training/distributed-systems/grpc/src/service/discovery"
	"github.com/kr/pretty"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	DEFAULT_ADDRESS = "localhost:50051"
)

func main() {
	address := flag.String("address", DEFAULT_ADDRESS, "address of server to connect to")
	uuid := flag.String("uuid", "", "uuid of node")
	uri := flag.String("uri", "", "uri of node")
	name := flag.String("name", "", "name of service the node belongs to")
	serviceType := flag.String("type", "", "type of service.  Either MasterMaster or MasterSlave")
	leader := flag.Bool("leader", false, "if the node is currently the leader")
	list := flag.String("list", "", "list nodes by service.  use an '*' for all nodes regardless of service")

	flag.Parse()

	if *list != "" {
		listNodes(*address, *list)
		return
	}

	// validate inputs
	if *uuid == "" {
		fmt.Println("uuid can not be blank")
		return
	}

	if *uri == "" {
		fmt.Println("uri can not be blank")
		return
	}

	if *name == "" {
		fmt.Println("name can not be blank")
		return
	}

	var st discovery.ServiceType
	if t, err := typeToEnum(*serviceType); err != nil {
		fmt.Printf("err: %s\n", err)
		return
	} else {
		st = t
	}

	request := discovery.RegistrationRequest{UUID: *uuid, Name: *name, Type: st, Leader: *leader, URI: *uri}
	register(*address, request)
}

func newClient(address string) (discovery.DiscoveryClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	c := discovery.NewDiscoveryClient(conn)
	return c, conn

}

func listNodes(address, name string) {
	// normalize out the * if they are asking to list all services
	if name == "*" {
		name = ""
	}
	c, conn := newClient(address)
	defer conn.Close()

	request := discovery.ListRequest{
		Name: name,
	}

	r, err := c.List(context.Background(), &request)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	pretty.Print(*r)
}

func register(address string, request discovery.RegistrationRequest) {
	c, conn := newClient(address)
	defer conn.Close()

	r, err := c.Register(context.Background(), &request)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	pretty.Print(*r)
}

func typeToEnum(s string) (discovery.ServiceType, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	switch s {
	case "mastermaster":
		return discovery.ServiceType_MasterMaster, nil
	case "masterslave":
		return discovery.ServiceType_MasterSlave, nil
	case "":
		return 0, fmt.Errorf("service type can not be blank")
	default:
		return 0, fmt.Errorf("%q is not a valid service type", s)
	}
}
