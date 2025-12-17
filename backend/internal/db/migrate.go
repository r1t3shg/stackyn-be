// Package db (migrate.go) handles database schema migrations.
// It uses Go's embed package to include SQL migration files in the binary.
package db

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

// migrationsFS is an embedded filesystem containing all SQL migration files.
// The //go:embed directive embeds all .sql files from the migrations directory
// into the binary at compile time, so migrations are always available.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate applies all pending database migrations in order.
// It tracks which migrations have been applied in a schema_migrations table.
// Migrations are applied in alphabetical order based on filename.
//
// Migration process:
// 1. Creates schema_migrations table if it doesn't exist
// 2. Reads all .sql files from the embedded migrations directory
// 3. Sorts them alphabetically to ensure consistent ordering
// 4. For each migration:
//   - Checks if it's already been applied
//   - If not, executes the SQL and records it in schema_migrations
//
// Returns:
//   - error: Any error that occurred during migration (table creation, file reading, SQL execution, etc.)
func (d *DB) Migrate() error {
	// Step 1: Create the schema_migrations tracking table if it doesn't exist.
	// This table stores which migrations have been applied to prevent duplicate execution.
	_, err := d.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Step 2: Read all files from the embedded migrations directory.
	// The migrationsFS is an embedded filesystem containing our SQL files.
	files, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	// Step 3: Filter and sort migration files.
	// We only process .sql files and sort them alphabetically to ensure consistent ordering.
	var migrationFiles []string
	for _, file := range files {
		// Only include files with .sql extension
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	// Sort alphabetically so migrations run in a predictable order
	sort.Strings(migrationFiles)

	// Step 4: Apply each migration that hasn't been applied yet.
	for _, filename := range migrationFiles {
		// Check if this migration has already been applied
		var exists bool
		err := d.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		// Skip if already applied
		if exists {
			continue
		}

		// Read the migration SQL file from the embedded filesystem
		// path := filepath.Join("migrations", filename)
		path := "migrations/" + filename
		content, err := migrationsFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		// Execute the migration SQL
		_, err = d.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record that this migration has been applied
		_, err = d.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		fmt.Printf("Applied migration: %s\n", filename)
	}

	return nil
}
