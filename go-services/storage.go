package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store interface {
	GetUser(string) (*User, error)
	CreateUser(*User) error
}

type PostgersStore struct {
	db *sql.DB
}

func NewPostgersStore() (*PostgersStore, error) {
	db, err := sql.Open("postgres", "dbname=postgres user=postgres password=postgres sslmode=disable")
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

func (s *PostgersStore) Init() error {
	err := s.CreateUserTable()
	if err != nil {
		return err
	}

	err = s.CreateUser(&User{Email: "bob@bob.bob", Password: "password"})
	return err
}

func (s *PostgersStore) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgersStore) CreateUser(user *User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`

	_, err := s.db.Exec(query, user.Email, user.Password)
	return err
}

func (s *PostgersStore) GetUser(email string) (*User, error) {
	query := `SELECT email, password FROM users WHERE email = $1`

	row := s.db.QueryRow(query, email)
	user := &User{}
	if err := row.Scan(&user.Email, &user.Password); err != nil {
		return nil, err
	}

	return user, nil
}
