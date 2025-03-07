package main

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "local.domain/user_proto"
)

func TestUserService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	mockAuth := NewMockAuthProvider(ctrl)

	service := NewUserService(mockRepo, mockAuth)

	ctx := context.Background()
	req := &pb.RegisterUserRequest{
		Login:       "testuser",
		Password:    "password",
		Email:       "test@example.com",
		FirstName:   "Test",
		LastName:    "User",
		BirthDate:   timestamppb.New(time.Now()),
		PhoneNumber: "1234567890",
	}

	mockAuth.EXPECT().
		HashPassword(req.Password).
		Return("hashed_password", nil)

	mockRepo.EXPECT().
		CreateUser(ctx, gomock.Any()).
		Return(1, nil)

	mockAuth.EXPECT().
		GenerateAccessToken(1).
		Return("access_token", nil)

	mockAuth.EXPECT().
		GenerateRefreshToken(1).
		Return("refresh_token", nil)

	resp, err := service.RegisterUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)
}

func TestUserService_AuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	mockAuth := NewMockAuthProvider(ctrl)

	service := NewUserService(mockRepo, mockAuth)

	ctx := context.Background()
	req := &pb.AuthenticateUserRequest{
		Login:    "testuser",
		Password: "password",
	}

	mockRepo.EXPECT().
		GetUserByLogin(ctx, req.Login).
		Return(&User{
			ID:       1,
			Login:    req.Login,
			Password: "hashed_password",
		}, nil)

	mockAuth.EXPECT().
		CheckPassword(req.Password, "hashed_password").
		Return(true)

	mockAuth.EXPECT().
		GenerateAccessToken(1).
		Return("access_token", nil)

	mockAuth.EXPECT().
		GenerateRefreshToken(1).
		Return("refresh_token", nil)

	resp, err := service.AuthenticateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)
}

func TestUserService_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	mockAuth := NewMockAuthProvider(ctrl)

	service := NewUserService(mockRepo, mockAuth)

	ctx := context.Background()
	req := &pb.GetProfileRequest{}

	mockAuth.EXPECT().
		GetTokenInfoFromContext(ctx).
		Return(&TokenInfo{
			UserID:    1,
			TokenType: AccessToken,
		}, nil)

	mockRepo.EXPECT().
		GetUserByID(ctx, 1).
		Return(&User{
			ID:          1,
			Login:       "testuser",
			Email:       "test@example.com",
			FirstName:   "Test",
			LastName:    "User",
			BirthDate:   time.Now(),
			PhoneNumber: "1234567890",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)

	resp, err := service.GetProfile(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", resp.Login)
	assert.Equal(t, "test@example.com", resp.Email)
}

func TestUserService_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	mockAuth := NewMockAuthProvider(ctrl)

	service := NewUserService(mockRepo, mockAuth)

	ctx := context.Background()
	req := &pb.UpdateProfileRequest{
		Email:       "updated@example.com",
		FirstName:   "Updated",
		LastName:    "User",
		BirthDate:   timestamppb.New(time.Now()),
		PhoneNumber: "0987654321",
	}

	mockAuth.EXPECT().
		GetTokenInfoFromContext(ctx).
		Return(&TokenInfo{
			UserID:    1,
			TokenType: AccessToken,
		}, nil)

	mockRepo.EXPECT().
		GetUserByID(ctx, 1).
		Return(&User{
			ID:          1,
			Login:       "testuser",
			Email:       "test@example.com",
			FirstName:   "Test",
			LastName:    "User",
			BirthDate:   time.Now(),
			PhoneNumber: "1234567890",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)

	mockRepo.EXPECT().
		UpdateUser(ctx, gomock.Any()).
		Return(nil)

	resp, err := service.UpdateProfile(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUserService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	mockAuth := NewMockAuthProvider(ctrl)

	service := NewUserService(mockRepo, mockAuth)

	ctx := context.Background()
	req := &pb.RefreshTokenRequest{
		RefreshToken: "refresh_token",
	}

	mockAuth.EXPECT().
		GetTokenInfoFromToken(req.RefreshToken).
		Return(&TokenInfo{
			UserID:    1,
			TokenType: RefreshToken,
		}, nil)

	mockAuth.EXPECT().
		GenerateAccessToken(1).
		Return("new_access_token", nil)

	resp, err := service.RefreshToken(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "new_access_token", resp.AccessToken)
}
