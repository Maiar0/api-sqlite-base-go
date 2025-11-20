package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Registers the SQLite driver with database/sql
)

// Store represents a base directory where game DB files are stored.
type Store struct{ BaseDir string }

// New creates a new Store bound to the given base directory.
// Ensures the directory exists; returns a pointer to Store.
func New(base string) *Store {
	_ = os.MkdirAll(base, 0o755) // TODO: add error handling
	return &Store{BaseDir: base}
}

// OpenFor opens (or creates) a SQLite DB file for the given game ID
// and ensures the schema is applied from the provided schema file.
func (s *Store) OpenFor(gameID string, schemaPath string) (*sql.DB, error) {
	path := filepath.Join(s.BaseDir, gameID+".db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := ensureSchema(db, schemaPath); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// ensureSchema reads the SQL schema file at the given path and executes it
// against the provided database connection.
func ensureSchema(db *sql.DB, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(b)) // SQLite accepts multi-statements separated by ';'
	return err
}
