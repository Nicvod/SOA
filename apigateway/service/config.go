package main

import (
	"flag"
	"fmt"
)

type Config struct {
	UserServiceEndpoint string
	PostServiceEndpoint string
	RestEndpoint        string
	PublicKeyFile       string
}

func NewConfig() (*Config, error) {
	var userServiceEndpoint, postServiceEndpoint, publicFile string
	var restEndpoint int
	flag.StringVar(&userServiceEndpoint, "user_service_endpoint", "user_app:50051", "service port")
	flag.StringVar(&postServiceEndpoint, "post_service_endpoint", "post_app:50051", "service port")
	flag.IntVar(&restEndpoint, "rest_endpoint", 80, "service port")
	flag.StringVar(&publicFile, "public_key", "", "path to JWT public key `file`")
	flag.Parse()
	if publicFile == "" {
		return nil, fmt.Errorf("public key file is not provided")
	}
	return &Config{
		UserServiceEndpoint: userServiceEndpoint,
		PostServiceEndpoint: postServiceEndpoint,
		RestEndpoint:        fmt.Sprint(restEndpoint),
		PublicKeyFile:       publicFile,
	}, nil
}
