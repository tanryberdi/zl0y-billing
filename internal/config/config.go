package config

import "os"

type Config struct {
	Port          string
	PostgresDSN   string
	MongoURI      string
	MongoDatabase string
	JWTSecret     string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		PostgresDSN:   getEnv("POSTGRES_DSN", "postgres://user:password@localhost:5432/billing?sslmode=disable"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "billing"),
		JWTSecret:     getEnv("JWT_SECRET", "my-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
