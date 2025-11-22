package auth

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	baseDir    = "auth"
	schemaPath = "auth/user.sql"
	dbFileName = "users.db"
)

type Store struct {
	db *sql.DB
}

var userStore *Store

func GetUserStore() (*Store, error) {
	if userStore == nil {
		err := InitUserDB()
		if err != nil {
			log.Printf("[store.go] Error getting user DB %e", err)
			return nil, err
		}
	}

	return userStore, nil
}

func CloseUserDB() error {
	if userStore == nil || userStore.db == nil {
		return nil
	}
	err := userStore.db.Close()
	userStore = nil
	log.Printf("[store.go] User DB closed successfully")
	return err
}

func InitUserDB() error {
	// Ensure baseDir exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("creating base dir %q: %w", baseDir, err)
	}

	dbPath := filepath.Join(baseDir, dbFileName)

	// Open (or create) the SQLite DB file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("opening db %q: %w", dbPath, err)
	}

	// If anything fails after this, close db before returning
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		db.Close()
		return fmt.Errorf("reading schema file %q: %w", schemaPath, err)
	}

	if _, err := db.Exec(string(schemaBytes)); err != nil {
		db.Close()
		return fmt.Errorf("executing schema from %q: %w", schemaPath, err)
	}

	userStore := &Store{db: db}
	_ = userStore // to avoid unused variable warning

	log.Printf("[store.go] User DB initialized successfully at %s", dbPath)

	return nil
}

func (u *Store) NewUser(user string, email string, password string) (sql.Result, error) {
	log.Printf("[store.go] NewUser called for user: %s, email: %s", user, email)
	now := time.Now().UTC().Format(time.RFC3339)
	userUUID := uuid.NewString()

	result, err := u.db.Exec(`
	INSERT INTO users (userUUID, username, email, password) 
	VALUES (?, ?, ?, ?, ?)
	`, userUUID, user, email, password, now)
	if err != nil {
		log.Printf("[store.go] Error inserting new user: %v", err)
	}
	insertID, _ := result.LastInsertId()
	log.Printf("[store.go] New user inserted with ID: %d", insertID)

	return result, err
}
