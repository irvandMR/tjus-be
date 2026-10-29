package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	DatabaseUrl string
	Port        string
	JwtSecret   string
}

func LoadEnvironment() *Environment {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Environment{
		DatabaseUrl: os.Getenv("DB_URL"),
		Port:        os.Getenv("PORT"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
	}
}
