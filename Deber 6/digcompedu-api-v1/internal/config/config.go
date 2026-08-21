// internal/config/config.go

package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	DBPath   string
	GroqKey  string
	GroqModel string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "digcompedu.db"),
		GroqKey:   os.Getenv("GROQ_API_KEY"),
		GroqModel: getEnv("GROQ_MODEL", "llama-3.3-70b-versatile"),
	}

	if cfg.GroqKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY es obligatorio")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
