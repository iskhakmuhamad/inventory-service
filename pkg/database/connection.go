package database

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var migrationFiles embed.FS

func Connect(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sqlx.DB) error {
	// Create migrations table if not exists
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	if _, err := db.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Read migration files from the embedded filesystem
	migrationEntries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Sort migration files
	var migrationNames []string
	for _, entry := range migrationEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationNames = append(migrationNames, entry.Name())
		}
	}
	sort.Strings(migrationNames)

	// Execute migrations
	for _, filename := range migrationNames {
		// Check if migration already executed
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM migrations WHERE filename = ?", filename)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count > 0 {
			continue // Migration already executed
		}

		// Read migration file from the embedded filesystem
		content, err := migrationFiles.ReadFile(fmt.Sprintf("migrations/%s", filename))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Execute migration
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record migration as executed
		if _, err := db.Exec("INSERT INTO migrations (filename) VALUES (?)", filename); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		fmt.Printf("Executed migration: %s\n", filename)
	}

	return nil
}
