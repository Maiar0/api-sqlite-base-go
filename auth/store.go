package auth

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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

func GetUserStore() (*Store, error) { //will initialize the store if not yet done
	if userStore == nil {
		log.Printf("[store.go] userStore == nil")
		store, err := InitUserDB()
		if err != nil {
			log.Printf("[store.go] Error getting user DB %v", err)
			return nil, err
		}
		userStore = store
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

func InitUserDB() (*Store, error) {
	// Ensure baseDir exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Printf("[store.go] Error creating base dir: %q", err.Error())
		return nil, fmt.Errorf("creating base dir %q: %w", baseDir, err)
	}

	dbPath := filepath.Join(baseDir, dbFileName)

	// Open (or create) the SQLite DB file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("[store.go] Error creating/opening db: %q", err.Error())
		return nil, fmt.Errorf("opening db %q: %w", dbPath, err)
	}

	// If anything fails after this, close db before returning
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("reading schema file %q: %w", schemaPath, err)
	}

	if _, err := db.Exec(string(schemaBytes)); err != nil {
		db.Close()
		return nil, fmt.Errorf("executing schema from %q: %w", schemaPath, err)
	}

	log.Printf("[store.go] User DB initialized successfully at %s", dbPath)

	return &Store{db: db}, nil
}

func (u *Store) NewUser(user string, email string, password string) (sql.Result, error) {
	log.Printf("[store.go] NewUser called for user: %s, email: %s", user, email)
	userUUID := uuid.NewString()
	password_hash, err := HashPassword(password) // TODO: hash the password before storing
	if err != nil {
		log.Printf("[store.go] Error hashing password: %v", err)
		return nil, err
	}
	result, err := u.db.Exec(`
	INSERT INTO users (uuid, username, email, password_hash) 
	VALUES (?, ?, ?, ?)
	`, userUUID, user, email, password_hash)
	if err != nil {
		log.Printf("[store.go] Error inserting new user: %v", err)
	}
	insertID, _ := result.LastInsertId()
	log.Printf("[store.go] New user inserted with ID: %d", insertID)

	return result, err
}

type Row struct {
	ID           int
	UUID         string
	username     string
	email        string
	passwordHash string
	created_at   string
	updated_at   string
	isActive     bool
	lastLogin    string
}

func (u *Store) GetUserByUsername(username string) (*Row, error) {
	log.Printf("[store.go] GetUserByUsername called for username: %s", username)
	row := u.db.QueryRow(`
	SELECT id, uuid, username, email, password_hash, created_at, updated_at, is_active, last_login_at 
	FROM users 
	WHERE username = ?
	`, username)

	var user Row
	err := row.Scan(
		&user.ID, &user.UUID, &user.username,
		&user.email, &user.passwordHash, &user.created_at,
		&user.updated_at, &user.isActive, &user.lastLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[store.go] No user found with username: %s", username)
			return nil, nil
		}
		log.Printf("[store.go] Error querying user by username: %v", err)
		return nil, err
	}
	log.Printf("[store.go] User found: %+v", user)
	return &user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
