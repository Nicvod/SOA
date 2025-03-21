package main

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

type TokenKeys struct {
	jwtPrivate *rsa.PrivateKey
	jwtPublic  *rsa.PublicKey
}

type TokenInfo struct {
	UserID    int
	TokenType TokenType
}

type TokenClaims struct {
	UserID    int       `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	jwtLib.RegisteredClaims
}

type AuthProvider interface {
	GenerateAccessToken(userID int) (string, error)
	GenerateRefreshToken(userID int) (string, error)
	GetTokenInfoFromToken(tokenString string) (*TokenInfo, error)
	GetAuthTokenFromContext(ctx context.Context) (string, error)
	GetTokenInfoFromContext(ctx context.Context) (*TokenInfo, error)
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) bool
}

type TokenManager struct {
	keys TokenKeys
}

func NewTokenManager(cfg *Config) (AuthProvider, error) {
	privateKeyData, err := os.ReadFile(cfg.PrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := jwtLib.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	publicKeyData, err := os.ReadFile(cfg.PublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKey, err := jwtLib.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	return &TokenManager{
		keys: TokenKeys{
			jwtPrivate: privateKey,
			jwtPublic:  publicKey,
		},
	}, nil
}

func (tm *TokenManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (tm *TokenManager) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (tm *TokenManager) GenerateAccessToken(userID int) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		TokenType: AccessToken,
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodRS256, claims)
	return token.SignedString(tm.keys.jwtPrivate)
}

func (tm *TokenManager) GenerateRefreshToken(userID int) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		TokenType: RefreshToken,
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodRS256, claims)
	return token.SignedString(tm.keys.jwtPrivate)
}

func (tm *TokenManager) GetTokenInfoFromToken(tokenString string) (*TokenInfo, error) {
	token, err := jwtLib.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwtLib.Token) (interface{}, error) {
		return tm.keys.jwtPublic, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return &TokenInfo{
			UserID:    claims.UserID,
			TokenType: claims.TokenType,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (tm *TokenManager) GetAuthTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is not provided")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", errors.New("authorization token is missing")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == "" {
		return "", errors.New("authorization token is empty")
	}

	return token, nil
}

func (tm *TokenManager) GetTokenInfoFromContext(ctx context.Context) (*TokenInfo, error) {
	authToken, err := tm.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return tm.GetTokenInfoFromToken(authToken)
}
