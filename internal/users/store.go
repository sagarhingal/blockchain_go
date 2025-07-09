package users

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, password TEXT)`); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO users(username, password) VALUES(?, ?)`, username, string(hash))
	return err
}

func (s *Store) getHash(username string) (string, error) {
	var hash string
	err := s.db.QueryRow(`SELECT password FROM users WHERE username=?`, username).Scan(&hash)
	return hash, err
}

func (s *Store) ValidateUser(username, password string) error {
	hash, err := s.getHash(username)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}

func (s *Store) UpdatePassword(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`UPDATE users SET password=? WHERE username=?`, string(hash), username)
	return err
}
