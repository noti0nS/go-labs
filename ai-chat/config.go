package main

import (
	"os"

	"github.com/joho/godotenv"
)

type APIConfig struct {
	Url   string
	Key   string
	Model string
}

func LoadConfig() (*APIConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &APIConfig{
		Url:   os.Getenv("API_URL"),
		Key:   os.Getenv("API_KEY"),
		Model: os.Getenv("API_MODEL")}, nil
}
