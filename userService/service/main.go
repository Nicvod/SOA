package main

import (
	"fmt"
	"log"
	"net"

	pb "local.domain/user_proto"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("failed to get congig: %v", err)
	}
	db, err := NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	tokenManager, err := NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("failed to create token manager: %v", err)
	}
	service := NewUserService(repo, tokenManager)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, service)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ServicePort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server is running on port %s\n", cfg.ServicePort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
