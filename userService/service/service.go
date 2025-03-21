package main

import (
	"context"
	"log"
	"time"

	pb "local.domain/user_proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	repo         UserRepository
	authProvider AuthProvider
	pb.UnimplementedUserServiceServer
}

func NewUserService(repo UserRepository, tokenManager AuthProvider) *UserService {
	return &UserService{repo: repo, authProvider: tokenManager}
}

func (s *UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	hashedPassword, err := s.authProvider.HashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	user := &User{
		Login:       req.Login,
		Password:    hashedPassword,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		BirthDate:   req.BirthDate.AsTime(),
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	accessToken, err := s.authProvider.GenerateAccessToken(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	refreshToken, err := s.authProvider.GenerateRefreshToken(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.RegisterUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	user, err := s.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if !s.authProvider.CheckPassword(req.Password, user.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	accessToken, err := s.authProvider.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	refreshToken, err := s.authProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.AuthenticateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	tokenInfo, err := s.authProvider.GetTokenInfoFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to get token: %v", err)
	}
	if tokenInfo.TokenType != AccessToken {
		return nil, status.Errorf(codes.PermissionDenied, "bad token type")
	}

	user, err := s.repo.GetUserByID(ctx, tokenInfo.UserID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.GetProfileResponse{
		Login:       user.Login,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		BirthDate:   timestamppb.New(user.BirthDate),
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	tokenInfo, err := s.authProvider.GetTokenInfoFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to get token: %v", err)
	}
	if tokenInfo.TokenType != AccessToken {
		return nil, status.Errorf(codes.PermissionDenied, "bad token type")
	}

	user, err := s.repo.GetUserByID(ctx, tokenInfo.UserID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.BirthDate = req.BirthDate.AsTime()
	user.PhoneNumber = req.PhoneNumber
	user.UpdatedAt = time.Now()

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UpdateProfileResponse{}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	tokenInfo, err := s.authProvider.GetTokenInfoFromToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to get info from token: %v", err)
	}
	log.Println(tokenInfo)
	if tokenInfo.TokenType != RefreshToken {
		return nil, status.Error(codes.InvalidArgument, "bad token type")
	}
	newAccessToken, err := s.authProvider.GenerateAccessToken(tokenInfo.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate new access token: %v", err)
	}
	return &pb.RefreshTokenResponse{AccessToken: newAccessToken}, nil
}
