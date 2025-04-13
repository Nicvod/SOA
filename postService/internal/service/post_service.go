package service

import (
	"context"
	"fmt"
	"log"

	postgres "github.com/Nicvod/SOA/postService/internal/repository"
	post_proto "github.com/Nicvod/SOA/postService/post_proto"
	"github.com/Nicvod/SOA/utils/auth"
)

type PostService struct {
	repo       *postgres.PostRepository
	authHelper auth.AuthProvider
}

func NewPostService(repo *postgres.PostRepository, authHelper auth.AuthProvider) *PostService {
	return &PostService{
		repo:       repo,
		authHelper: authHelper,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *post_proto.CreatePostRequest) (*post_proto.PostResponse, error) {
	tokenInfo, err := s.authHelper.TokenInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.CreatePost(ctx, req, fmt.Sprint(tokenInfo.UserID))
}

func (s *PostService) GetPost(ctx context.Context, req *post_proto.GetPostRequest) (*post_proto.PostResponse, error) {
	tokenInfo, err := s.authHelper.TokenInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.GetPost(ctx, req.PostId, fmt.Sprint(tokenInfo.UserID))
}

func (s *PostService) UpdatePost(ctx context.Context, req *post_proto.UpdatePostRequest) (*post_proto.PostResponse, error) {
	tokenInfo, err := s.authHelper.TokenInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdatePost(ctx, req, fmt.Sprint(tokenInfo.UserID))
}

func (s *PostService) DeletePost(ctx context.Context, req *post_proto.DeletePostRequest) error {
	tokenInfo, err := s.authHelper.TokenInfoFromContext(ctx)
	if err != nil {
		return err
	}
	log.Println(tokenInfo, req.PostId)

	return s.repo.DeletePost(ctx, req.PostId, fmt.Sprint(tokenInfo.UserID))
}

func (s *PostService) ListPosts(ctx context.Context, req *post_proto.ListPostsRequest) (*post_proto.ListPostsResponse, error) {
	tokenInfo, err := s.authHelper.TokenInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.ListPosts(ctx, fmt.Sprint(tokenInfo.UserID), req.Page, req.PageSize)
}
