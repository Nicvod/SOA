package models

import "errors"

var (
	ErrPostNotFound = errors.New("post not found")
	ErrUnauthorized = errors.New("unauthorized access")
)
