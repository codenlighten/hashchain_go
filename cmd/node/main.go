package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/greg/hashchain/pkg/node"
	pb "github.com/greg/hashchain/proto"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	nodeService := node.NewNodeService()
	pb.RegisterNodeServiceServer(s, nodeService)

	log.Printf("Starting node server on port %d", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
