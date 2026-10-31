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
	AppEnv      string
}

func LoadEnvironment() *Environment {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env tidak ditemukan, membaca dari environment sistem")
	}

	return &Environment{
		DatabaseUrl: os.Getenv("DB_URL"),
		Port:        os.Getenv("PORT"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
		AppEnv:      os.Getenv("APP_ENV"),
	}
}

func (e Environment) IsProduction() bool {
	return e.AppEnv == "production"
}
