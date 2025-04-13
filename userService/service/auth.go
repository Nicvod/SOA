package main

import (
	"fmt"
	"os"

	auth "github.com/Nicvod/SOA/utils/auth"
	jwtLib "github.com/golang-jwt/jwt/v5"
)

func NewTokenManager(cfg *Config) (*auth.TokenKeys, error) {
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

	return auth.NewTokenKeys(privateKey, publicKey), nil
}
