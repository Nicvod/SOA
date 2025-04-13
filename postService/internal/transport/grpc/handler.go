package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Nicvod/SOA/postService/internal/service"
	post_proto "github.com/Nicvod/SOA/postService/post_proto"
)

type PostHandler struct {
	post_proto.UnimplementedPostServiceServer
	service *service.PostService
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

func (h *PostHandler) CreatePost(ctx context.Context, req *post_proto.CreatePostRequest) (*post_proto.PostResponse, error) {
	return h.service.CreatePost(ctx, req)
}

func (h *PostHandler) GetPost(ctx context.Context, req *post_proto.GetPostRequest) (*post_proto.PostResponse, error) {
	return h.service.GetPost(ctx, req)
}

func (h *PostHandler) UpdatePost(ctx context.Context, req *post_proto.UpdatePostRequest) (*post_proto.PostResponse, error) {
	return h.service.UpdatePost(ctx, req)
}

func (h *PostHandler) DeletePost(ctx context.Context, req *post_proto.DeletePostRequest) (*emptypb.Empty, error) {
	err := h.service.DeletePost(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *PostHandler) ListPosts(ctx context.Context, req *post_proto.ListPostsRequest) (*post_proto.ListPostsResponse, error) {
	return h.service.ListPosts(ctx, req)
}
