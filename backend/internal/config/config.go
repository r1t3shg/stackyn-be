// Package config provides configuration management for the application.
// It loads configuration from environment variables with sensible defaults.
package config

import (
	"os"
)

// Config holds all application configuration values.
// These values are loaded from environment variables at startup.
type Config struct {
	// DatabaseURL is the PostgreSQL connection string.
	// Format: postgres://user:password@host:port/database?sslmode=disable
	// Default: postgres://postgres:postgres@localhost:5432/mvp?sslmode=disable
	DatabaseURL string

	// DockerHost is the address of the Docker daemon.
	// Can be a Unix socket (unix:///var/run/docker.sock) or TCP address (tcp://host:port).
	// Default: unix:///var/run/docker.sock
	DockerHost string

	// BaseDomain is the base domain used for subdomain routing.
	// Deployed apps will be accessible at {subdomain}.{BaseDomain}
	// Default: localhost
	BaseDomain string

	// Port is the port number for the HTTP API server.
	// Default: 8080
	Port string
}

// Load reads configuration from environment variables and returns a Config struct.
// If an environment variable is not set, it uses the provided default value.
// This function should be called at application startup.
//
// Returns:
//   - *Config: A pointer to a Config struct with all values populated
func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:ritesh@localhost:5432/mvp?sslmode=disable"),
		DockerHost:  getEnv("DOCKER_HOST", "tcp://localhost:2375"),
		BaseDomain:  getEnv("BASE_DOMAIN", "localhost"),
		Port:        getEnv("PORT", "8080"),
	}
}

// getEnv retrieves an environment variable value, returning the default if not set.
// This is a helper function used internally by Load().
//
// Parameters:
//   - key: The name of the environment variable to read
//   - defaultValue: The value to return if the environment variable is not set or empty
//
// Returns:
//   - string: The environment variable value, or defaultValue if not set
func getEnv(key, defaultValue string) string {
	// Try to get the environment variable
	if value := os.Getenv(key); value != "" {
		return value
	}
	// Return default if not set or empty
	return defaultValue
}
