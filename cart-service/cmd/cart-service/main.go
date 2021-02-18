package main

import (
	"cart-service/internal/auth_service"
	"context"
	"log"

	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	client := auth_service.NewAuthClient(conn)

	response, err := client.Authenticate(context.Background(), &auth_service.AuthRequest{Token: "123456"})
	if err != nil {
		log.Fatalf("Error when calling Authenticate: %v", err)
	}
	log.Printf("Response from server: %b", response.Authenticated)
}
