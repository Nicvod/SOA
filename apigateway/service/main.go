package main

import (
	"log"

	post_proto "github.com/Nicvod/SOA/postService/post_proto"
	user_proto "github.com/Nicvod/SOA/userService/user_proto"

	"github.com/Nicvod/SOA/utils/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userClient   user_proto.UserServiceClient
	postClient   post_proto.PostServiceClient
	authProvider auth.AuthProvider
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("bad config: %v", err)
	}
	authProvider, err = NewAuthProvider(cfg)
	if err != nil {
		log.Fatalf("failed to create auth provider: %v", err)
	}
	log.Println(cfg.UserServiceEndpoint)
	userConn, err := grpc.NewClient(cfg.UserServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer userConn.Close()

	postConn, err := grpc.NewClient(cfg.PostServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer postConn.Close()

	userClient = user_proto.NewUserServiceClient(userConn)
	postClient = post_proto.NewPostServiceClient(postConn)

	r := gin.Default()

	r.Use(MiddlewareWrapper())

	api := r.Group("/api")
	{
		api.Static("/swagger", "/app/apigateway/service/swagger")
		api.POST("/v1/register", registerUser)
		api.POST("/v1/authenticate", authenticateUser)
		api.POST("/v1/refresh-token", refreshToken)
		api.PUT("/v1/profile", updateProfile)
		api.GET("/v1/profile", getProfile)
		posts := api.Group("/v1/posts")
		{
			posts.POST("", createPost)
			posts.GET("", listPosts)

			postID := posts.Group("/:post_id")
			{
				postID.GET("", getPost)
				postID.PUT("", updatePost)
				postID.DELETE("", deletePost)
			}
		}
	}

	if err := r.Run(":" + cfg.RestEndpoint); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
