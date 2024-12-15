package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
	QueueHost  string
	QueuePort  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		RedisHost:  os.Getenv("REDIS_HOST"),
		RedisPort:  os.Getenv("REDIS_PORT"),
		QueueHost:  os.Getenv("QUEUE_HOST"),
		QueuePort:  os.Getenv("QUEUE_PORT"),
	}

	return config, nil
}
