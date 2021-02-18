package main

import (
	"context"
	"log"
	"net"

	auth_service "auth-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	auth_service.UnimplementedAuthServer
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen to port 9000: %v", err)
	}

	server := server{}

	srv := grpc.NewServer()
	auth_service.RegisterAuthServer(srv, &server)
	reflection.Register(srv)

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("gRPC failed to serve on port 9000: %v", err)
	}
}

// Auth method
func (s *server) Authenticate(ctx context.Context, request *auth_service.AuthRequest) (*auth_service.AuthResponse, error) {
	// The JWT string in the request
	jwtString := request.GetToken()

	// Fake verify the string
	verified := len(jwtString) > 5

	// Return whether the jwt is verified
	return &auth_service.AuthResponse{Authenticated: verified}, nil
}
