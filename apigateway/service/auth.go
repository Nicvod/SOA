package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	auth "github.com/Nicvod/SOA/utils/auth"
	"github.com/gin-gonic/gin"
	jwtLib "github.com/golang-jwt/jwt/v5"
)

func NewAuthProvider(cfg *Config) (*auth.TokenKeys, error) {
	publicKeyData, err := os.ReadFile(cfg.PublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKey, err := jwtLib.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	return auth.NewTokenKeys(nil, publicKey), nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		tokenInfo, err := authProvider.ValidateToken(token, auth.AccessToken)
		if err != nil {
			log.Println("Bad auth:", token, tokenInfo, err)
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

var (
	noAuthPaths = []string{
		"/api/v1/register",
		"/api/v1/authenticate",
		"/api/v1/refresh-token",
		"/api/swagger",
	}
)

func skipAuth(path string) bool {
	skip := false
	for _, skipPath := range noAuthPaths {
		if strings.HasPrefix(path, skipPath) {
			skip = true
		}
	}
	return skip
}

func MiddlewareWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipAuth(c.Request.URL.Path) {
			c.Next()
			return
		}
		AuthMiddleware()(c)
	}
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
