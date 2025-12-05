package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"pubsub-grpc/cmd/server/internal/service/pubsub"
	pb "pubsub-grpc/pkg/pb/pubsub"

	"google.golang.org/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "port", "12345", "The port to listen on")
	flag.Parse()
}

func main() {
	server := grpc.NewServer()

	pb.RegisterPubsubServiceServer(server, &pubsub.PubsubService{})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
