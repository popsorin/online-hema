package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"hema-lessons/internal/database"
)

type Config = database.Config

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	cfg := SetupTestConfig()
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	CleanDatabase(t, db)

	if err := database.RunMigrations(db, "../../migrations/test"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

func CleanDatabase(t *testing.T, db *sql.DB) {
	t.Helper()

	tables := []string{
		"subscriptions",
		"users",
		"techniques",
		"chapters",
		"fighting_books",
		"sword_masters",
		"schema_migrations",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			t.Logf("warning: failed to drop table %s: %v", table, err)
		}
	}

	_, _ = db.Exec("DROP TYPE IF EXISTS subscription_status CASCADE")
}

func TeardownTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Errorf("failed to close database: %v", err)
	}
}

func SetupTestConfig() *database.Config {
	return &database.Config{
		Host:     getRequiredEnv("TEST_DATABASE_HOST"),
		Port:     getRequiredEnv("TEST_DATABASE_PORT"),
		User:     getRequiredEnv("TEST_DATABASE_USER"),
		Password: getRequiredEnv("TEST_DATABASE_PASSWORD"),
		DBName:   getRequiredEnv("TEST_DATABASE_DBNAME"),
		SSLMode:  getEnvOrDefault("TEST_DATABASE_SSLMODE", "disable"),
	}
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s not set", key))
	}
	return value
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
