package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

// User represents a user in the system
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Store represents a data store for the auth service
type Store interface {
	GetUser(string) (*User, error)
	CreateUser(*User) error
}

// PostgersStore represents a PostgreSQL data store
type PostgersStore struct {
	db *sql.DB
}

// NewPostgersStore creates a new PostgerSQL Store instance
func NewPostgersStore() (*PostgersStore, error) {
	conStr := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres")
	if err != nil {
		return nil, err
	}

	store := &PostgersStore{db: db}

	err = store.Init()
	if err != nil {
		return nil, err
	}

	return store, nil
}

// Init initializes the PostgersStore instance
func (s *PostgersStore) Init() error {
	err := s.CreateUserTable()
	if err != nil {
		return err
	}

	err = s.CreateUser(&User{Email: "bob@bob.bob", Password: "password"})
	return err
}

// CreateUserTable creates the user table in the database
func (s *PostgersStore) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL)`

	_, err := s.db.Exec(query)
	return err
}

// CreateUser creates a new user in the database
func (s *PostgersStore) CreateUser(user *User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`

	_, err := s.db.Exec(query, user.Email, user.Password)
	return err
}

// GetUser retrieves a user from the database
func (s *PostgersStore) GetUser(email string) (*User, error) {
	query := `SELECT email, password FROM users WHERE email = $1`

	row := s.db.QueryRow(query, email)
	user := &User{}
	if err := row.Scan(&user.Email, &user.Password); err != nil {
		return nil, err
	}

	return user, nil
}
