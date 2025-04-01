package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Nicvod/SOA/postService/internal/service"
	post_proto "github.com/Nicvod/SOA/postService/post_proto"

	"github.com/Nicvod/SOA/postService/internal/config"
	postgres "github.com/Nicvod/SOA/postService/internal/repository"
	"github.com/Nicvod/SOA/utils/auth"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	grpcServer *grpc.Server
	config     *config.Config
}

func NewServer(cfg *config.Config, db *sqlx.DB, authHelper auth.AuthProvider) *Server {
	postRepo := postgres.NewPostRepository(db)
	postService := service.NewPostService(postRepo, authHelper)
	postHandler := NewPostHandler(postService)

	grpcServer := grpc.NewServer()
	post_proto.RegisterPostServiceServer(grpcServer, postHandler)

	return &Server{
		grpcServer: grpcServer,
		config:     cfg,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", ":"+s.config.ServicePort)
	if err != nil {
		return err
	}

	log.Printf("Starting gRPC server on port %s", s.config.ServicePort)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
