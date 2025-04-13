package auth

import (
	"fmt"
	"os"

	cfg "github.com/Nicvod/SOA/postService/internal/config"
	authUtils "github.com/Nicvod/SOA/utils/auth"
	jwtLib "github.com/golang-jwt/jwt/v5"
)

func NewAuthProvider(cfg *cfg.Config) (*authUtils.TokenKeys, error) {
	publicKeyData, err := os.ReadFile(cfg.PublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKey, err := jwtLib.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	return authUtils.NewTokenKeys(nil, publicKey), nil
}
