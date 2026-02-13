package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	files, err := getMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	for _, file := range files {
		if err := executeMigration(db, filepath.Join(migrationsPath, file)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := db.Exec(query)
	return err
}

func getMigrationFiles(migrationsPath string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".up.sql") {
			files = append(files, d.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func executeMigration(db *sql.DB, filepath string) error {
	version := getMigrationVersion(filepath)
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Applied migration: %s\n", version)
	return nil
}

func getMigrationVersion(filepath string) string {
	base := strings.TrimSuffix(filepath, ".up.sql")
	return strings.TrimPrefix(base, "migrations/")
}
