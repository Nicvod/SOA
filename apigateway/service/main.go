package main

import (
	"log"

	user_proto "local.domain/user_proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcClient user_proto.UserServiceClient
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("bad config: %v", err)
	}
	log.Println("user_app:" + cfg.UserServiceEndpoint)
	conn, err := grpc.NewClient("user_app:"+cfg.UserServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	grpcClient = user_proto.NewUserServiceClient(conn)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.Static("/swagger", "./swagger")
		api.POST("/v1/register", registerUser)
		api.POST("/v1/authenticate", authenticateUser)
		api.POST("/v1/refresh-token", refreshToken)
		api.PUT("/v1/profile", updateProfile)
		api.GET("/v1/profile", getProfile)
	}

	if err := r.Run(":" + cfg.RestEndpoint); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
