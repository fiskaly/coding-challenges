package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	ListenAddress string
	Port          string

	// Crypto configuration
	RSAKeySize int
	ECCCurve   string

	// Logging
	LogLevel string
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	return &Config{
		ListenAddress: getEnv("LISTEN_ADDRESS", "0.0.0.0"),
		Port:          getEnv("PORT", "8080"),
		RSAKeySize:    getEnvAsInt("RSA_KEY_SIZE", 512), // 512 for demo, use 2048+ in production
		ECCCurve:      getEnv("ECC_CURVE", "P384"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

// GetListenAddr returns the full listen address (host:port)
func (c *Config) GetListenAddr() string {
	return c.ListenAddress + ":" + c.Port
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

