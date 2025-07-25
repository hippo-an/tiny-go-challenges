package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tpc", port)

	if err != nil {
		log.Fatalf("failed to listen on port %s", port)
	}

	gs := grpc.NewServer()

	log.Printf("start gRPC server on port %s", port)

}
