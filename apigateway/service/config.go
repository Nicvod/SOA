package main

import (
	"flag"
	"fmt"
)

type Config struct {
	UserServiceEndpoint string
	RestEndpoint        string
}

func NewConfig() (*Config, error) {
	userServiceEndpoint := flag.Int("user_service_endpoint", 50051, "service port")
	restEndpoint := flag.Int("rest_endpoint", 80, "service port")
	flag.Parse()
	return &Config{
		UserServiceEndpoint: fmt.Sprint(*userServiceEndpoint),
		RestEndpoint:        fmt.Sprint(*restEndpoint),
	}, nil
}
