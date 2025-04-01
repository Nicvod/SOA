package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"

	post_proto "github.com/Nicvod/SOA/postService/post_proto"
)

func createPost(c *gin.Context) {
	var request struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		IsPrivate   bool     `json:"is_private"`
		Tags        []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer "+token),
	)

	grpcReq := &post_proto.CreatePostRequest{
		Title:       request.Title,
		Description: request.Description,
		IsPrivate:   request.IsPrivate,
		Tags:        request.Tags,
	}

	resp, err := postClient.CreatePost(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          resp.Id,
		"title":       resp.Title,
		"description": resp.Description,
		"is_private":  resp.IsPrivate,
		"tags":        resp.Tags,
		"created_at":  resp.CreatedAt.AsTime(),
	})
}

func listPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer "+token),
	)

	grpcReq := &post_proto.ListPostsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	resp, err := postClient.ListPosts(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var posts []gin.H
	for _, p := range resp.Posts {
		posts = append(posts, gin.H{
			"id":          p.Id,
			"title":       p.Title,
			"description": p.Description,
			"is_private":  p.IsPrivate,
			"tags":        p.Tags,
			"created_at":  p.CreatedAt.AsTime(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":       posts,
		"total_count": resp.TotalCount,
		"page":        resp.Page,
		"page_size":   resp.PageSize,
	})
}

func getPost(c *gin.Context) {
	postID := c.Param("post_id")

	grpcReq := &post_proto.GetPostRequest{
		PostId: postID,
	}

	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer "+token),
	)

	resp, err := postClient.GetPost(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          resp.Id,
		"title":       resp.Title,
		"description": resp.Description,
		"creator_id":  resp.CreatorId,
		"is_private":  resp.IsPrivate,
		"tags":        resp.Tags,
		"created_at":  resp.CreatedAt.AsTime(),
		"updated_at":  resp.UpdatedAt.AsTime(),
	})
}

func updatePost(c *gin.Context) {
	postID := c.Param("post_id")

	var request struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		IsPrivate   bool     `json:"is_private"`
		Tags        []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer "+token),
	)

	grpcReq := &post_proto.UpdatePostRequest{
		PostId:      postID,
		Title:       request.Title,
		Description: request.Description,
		IsPrivate:   request.IsPrivate,
		Tags:        request.Tags,
	}

	resp, err := postClient.UpdatePost(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          resp.Id,
		"title":       resp.Title,
		"description": resp.Description,
		"is_private":  resp.IsPrivate,
		"tags":        resp.Tags,
		"updated_at":  resp.UpdatedAt.AsTime(),
	})
}

func deletePost(c *gin.Context) {
	postID := c.Param("post_id")

	grpcReq := &post_proto.DeletePostRequest{
		PostId: postID,
	}

	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer "+token),
	)

	_, err := postClient.DeletePost(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
