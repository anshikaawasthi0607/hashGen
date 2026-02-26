package config

import "os"

type Config struct {
	ServerPort string // HTTP server port 
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}