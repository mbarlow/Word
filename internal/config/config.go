package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Env       string
	Port      int
	DBPath    string
	LogLevel  string
	LogFormat string
}

func Load() *Config {
	return &Config{
		Env:       getEnv("WORD_ENV", "development"),
		Port:      getEnvInt("WORD_PORT", 8080),
		DBPath:    getEnv("WORD_DB_PATH", "./data/word.db"),
		LogLevel:  getEnv("WORD_LOG_LEVEL", "info"),
		LogFormat: getEnv("WORD_LOG_FORMAT", "json"),
	}
}

func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
