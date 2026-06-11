package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string
	MongoDBName    string
	Port           string
	YoutubeAPIKey  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	return &Config{
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:   getEnv("MONGO_DB_NAME", "vicio_api"),
		Port:          getEnv("PORT", "8080"),
		YoutubeAPIKey: getEnv("YOUTUBE_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}