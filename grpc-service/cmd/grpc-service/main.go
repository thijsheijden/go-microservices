package main

import (
	"context"
	"grpc-service/internal/greeter"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	greeter.UnimplementedGreetingServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen to port 9000: %v", err)
	}

	server := server{}

	grpcServer := grpc.NewServer()
	greeter.RegisterGreetingServiceServer(grpcServer, &server)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC failed to serve on port 9000: %v", err)
	}
}

// Implementation of the Greet method
func (s *server) Greeter(ctx context.Context, request *greeter.Greeting) (*greeter.Response, error) {
	// Get this pods name
	podName := os.Getenv("POD_NAME")

	// Send back this pod's name
	return &greeter.Response{Response: podName}, nil
}
