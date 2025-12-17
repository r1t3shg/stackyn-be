// Package db provides database connection and management functionality.
// It wraps the standard database/sql package with PostgreSQL-specific features.
package db

import (
	"database/sql"
	"fmt"
	"log"

	// Import PostgreSQL driver (blank import for side effects)
	_ "github.com/lib/pq"
)

// DB wraps the standard sql.DB connection with additional methods.
// It embeds *sql.DB, so all standard database/sql methods are available.
type DB struct {
	*sql.DB
}

// New creates a new database connection using the provided connection string.
// It opens the connection, verifies it with a ping, and returns a DB wrapper.
//
// Parameters:
//   - databaseURL: PostgreSQL connection string (e.g., postgres://user:pass@host:port/db?sslmode=disable)
//
// Returns:
//   - *DB: A database connection wrapper, or nil on error
//   - error: Any error that occurred during connection (connection failure, ping failure, etc.)
func New(databaseURL string) (*DB, error) {
	// Open database connection using PostgreSQL driver
	// This doesn't actually establish a connection yet, just prepares it
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify the connection works by pinging the database
	// This actually establishes a connection and tests it
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	// Wrap the sql.DB in our custom DB struct
	return &DB{db}, nil
}

// Close closes the database connection.
// This should be called when the application shuts down to clean up resources.
//
// Returns:
//   - error: Any error that occurred during connection closure
func (d *DB) Close() error {
	return d.DB.Close()
}

