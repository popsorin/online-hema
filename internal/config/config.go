package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
	App    AppConfig
}

type ServerConfig struct {
	Addr              string
	ReadHeaderTimeout int
}

type AppConfig struct {
	Environment string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Addr:              getEnv("SERVER_ADDR", ":8080"),
			ReadHeaderTimeout: getEnvAsInt("SERVER_READ_HEADER_TIMEOUT", 5),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENVIRONMENT", "development"),
		},
	}

	if err := validate(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func validate(config *Config) error {
	if config.Server.Addr == "" {
		return fmt.Errorf("server address is required")
	}
	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
