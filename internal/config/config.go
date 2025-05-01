package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DB   string
}

// Load reads .env (if present) and returns a Config object.
func Load() Config {
	_ = godotenv.Load() // ignore error; .env is optional

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db := os.Getenv("DATABASE_URL")
	if db == "" {
		db = "healthcare-api.db"
	}
	return Config{Port: port, DB: db}
}

func MustLoad() Config {
	c := Load()
	if c.DB == "" {
		log.Fatal("database config not found")
	}
	return c
}
