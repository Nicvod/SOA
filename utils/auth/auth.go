package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidTokenType  = errors.New("invalid token type")
	ErrInvalidSigningKey = errors.New("invalid signing key")
	ErrNoTokenInContext  = errors.New("no token in context")
	ErrNoAuthInContext   = errors.New("no auth identity in context")
)

type AuthProvider interface {
	GenerateToken(info TokenInfo, expiresIn time.Duration) (string, error)
	ValidateToken(tokenString string, expectedType TokenType) (*TokenInfo, error)
	GetTokenInfo(tokenString string) (*TokenInfo, error)
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) bool
	TokenInfoFromContext(ctx context.Context) (*TokenInfo, error)
	TokenFromGRPCContext(ctx context.Context) (string, error)
	GRPCContextWithToken(ctx context.Context, token string) context.Context
}

type TokenKeys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewTokenKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *TokenKeys {
	return &TokenKeys{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

type TokenInfo struct {
	UserID    int
	UserLogin string
	TokenType TokenType
}

type TokenClaims struct {
	UserID    int       `json:"user_id"`
	UserLogin string    `json:"user_login"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

func (tk *TokenKeys) GenerateToken(info TokenInfo, expiresIn time.Duration) (string, error) {
	if tk.privateKey == nil {
		return "", ErrInvalidSigningKey
	}

	claims := TokenClaims{
		UserID:    info.UserID,
		UserLogin: info.UserLogin,
		TokenType: info.TokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(tk.privateKey)
}

func (tk *TokenKeys) ValidateToken(tokenString string, expectedType TokenType) (*TokenInfo, error) {
	if tk.publicKey == nil {
		return nil, ErrInvalidSigningKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tk.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != expectedType {
		return nil, ErrInvalidTokenType
	}

	return &TokenInfo{
		UserID:    claims.UserID,
		UserLogin: claims.UserLogin,
		TokenType: claims.TokenType,
	}, nil
}

func (tk *TokenKeys) GetTokenInfo(tokenString string) (*TokenInfo, error) {
	if tk.publicKey == nil {
		return nil, ErrInvalidSigningKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tk.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &TokenInfo{
		UserID:    claims.UserID,
		UserLogin: claims.UserLogin,
		TokenType: claims.TokenType,
	}, nil
}

const (
	authorizationHeader = "authorization"
	tokenPrefix         = "Bearer "
)

func (tk *TokenKeys) TokenInfoFromContext(ctx context.Context) (*TokenInfo, error) {
	token, err := tk.TokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}
	return tk.GetTokenInfo(token)
}

func (tk *TokenKeys) TokenFromGRPCContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrNoTokenInContext
	}

	authHeaders := md.Get(authorizationHeader)
	if len(authHeaders) == 0 {
		return "", ErrNoTokenInContext
	}

	token := authHeaders[0]
	if len(token) <= len(tokenPrefix) || token[:len(tokenPrefix)] != tokenPrefix {
		return "", ErrInvalidToken
	}

	return token[len(tokenPrefix):], nil
}

func (tk *TokenKeys) GRPCContextWithToken(ctx context.Context, token string) context.Context {
	md := metadata.Pairs(authorizationHeader, tokenPrefix+token)
	return metadata.NewOutgoingContext(ctx, md)
}

func (tk *TokenKeys) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (tk *TokenKeys) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
