package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	App      AppConfig
}

type ServerConfig struct {
	Addr              string
	ReadHeaderTimeout int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	URL string
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
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnv("DATABASE_PORT", "5432"),
			User:     getEnv("DATABASE_USER", "postgres"),
			Password: getEnv("DATABASE_PASSWORD", "postgres"),
			DBName:   getEnv("DATABASE_DBNAME", "hema_lessons"),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379/0"),
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
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
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
