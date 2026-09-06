package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	SessionSecret string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "user=web password=web dbname=web host=localhost port=5433 sslmode=disable"),
		SessionSecret: getEnv("SESSION_SECRET", "development-secret-blah-blah"),
	}
} 

func getEnv(key, defaultvalue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultvalue
	}
	return  value
}