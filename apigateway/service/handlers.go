package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	user_proto "local.domain/user_proto"
)

type CustomTimestamp struct {
	*timestamppb.Timestamp
}

func (ct *CustomTimestamp) UnmarshalJSON(data []byte) error {
	dateString := strings.Trim(string(data), `"`)
	t, err := time.Parse(time.RFC3339Nano, dateString)
	if err != nil {
		return err
	}
	ct.Timestamp = timestamppb.New(t)
	return nil
}

func (ct *CustomTimestamp) MarshalJSON() ([]byte, error) {
	if ct.Timestamp == nil {
		return []byte("null"), nil
	}
	t := ct.Timestamp.AsTime()
	return []byte(`"` + t.Format(time.RFC3339Nano) + `"`), nil
}

type RegisterUserRequest struct {
	Login       string           `json:"login"`
	Password    string           `json:"password"`
	Email       string           `json:"email"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	BirthDate   *CustomTimestamp `json:"birth_date"`
	PhoneNumber string           `json:"phone_number"`
}

type UpdateProfileRequest struct {
	Email       string           `json:"email"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	BirthDate   *CustomTimestamp `json:"birth_date"`
	PhoneNumber string           `json:"phone_number"`
}

type GetProfileResponse struct {
	Login       string           `json:"login"`
	Email       string           `json:"email"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	BirthDate   *CustomTimestamp `json:"birth_date"`
	PhoneNumber string           `json:"phone_number"`
	CreatedAt   *CustomTimestamp `json:"created_at"`
	UpdatedAt   *CustomTimestamp `json:"updated_at"`
}

func registerUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &user_proto.RegisterUserRequest{
		Login:       req.Login,
		Password:    req.Password,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		BirthDate:   req.BirthDate.Timestamp,
		PhoneNumber: req.PhoneNumber,
	}

	res, err := grpcClient.RegisterUser(context.Background(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func updateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	grpcReq := &user_proto.UpdateProfileRequest{
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		BirthDate:   req.BirthDate.Timestamp,
		PhoneNumber: req.PhoneNumber,
	}

	res, err := grpcClient.UpdateProfile(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func getProfile(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer "+token),
	)

	res, err := grpcClient.GetProfile(ctx, &user_proto.GetProfileRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profileResponse := GetProfileResponse{
		Login:       res.Login,
		Email:       res.Email,
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		BirthDate:   &CustomTimestamp{res.BirthDate},
		PhoneNumber: res.PhoneNumber,
		CreatedAt:   &CustomTimestamp{res.CreatedAt},
		UpdatedAt:   &CustomTimestamp{res.UpdatedAt},
	}

	c.JSON(http.StatusOK, profileResponse)
}

func authenticateUser(c *gin.Context) {
	var req user_proto.AuthenticateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := grpcClient.AuthenticateUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login or password"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func refreshToken(c *gin.Context) {
	var req user_proto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := grpcClient.RefreshToken(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
