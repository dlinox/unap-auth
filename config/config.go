package config

import (
	"log"
)

type Config struct {
	DBDSN     string // Data Source Name for the database connection
	JWTSecret string // Secret key for JWT
}

func LoadConfig() (*Config, error) {
	dsn := "root:@tcp(localhost:3306)/unap.core"
	if dsn == "" {
		log.Fatal("DATABASE_DSN environment variable is not set")
	}

	return &Config{
		DBDSN:     dsn,
		JWTSecret: "P9HDdrq4LZJ9RqFT6A7sC5kJxQPYypra8fLgteKWJhPmcFjaaBeRHrMRkRWX4y8V",
	}, nil
}
