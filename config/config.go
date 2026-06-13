package config

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	ServerName string
	Version    string
	LogLevel   slog.Level
}

func Load() *Config {
	return &Config{
		ServerName: getEnv("SERVER_NAME", "atlassian-mcp-extensions"),
		Version:    getEnv("VERSION", "v0.1.0"),
		LogLevel:   parseLogLevel(getEnv("LOG_LEVEL", "info")),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseLogLevel(value string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
