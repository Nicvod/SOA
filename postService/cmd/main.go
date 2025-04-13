package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	authInternal "github.com/Nicvod/SOA/postService/internal/auth"
	"github.com/Nicvod/SOA/postService/internal/config"
	"github.com/Nicvod/SOA/postService/internal/transport/grpc"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	authHelper, err := authInternal.NewAuthProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to create authhelper: %v", err)
	}
	server := grpc.NewServer(cfg, db, authHelper)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.Stop()
	log.Println("Server stopped gracefully")
}
