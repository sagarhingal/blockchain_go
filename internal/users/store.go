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

// User represents an actor in the system.
// Email is used as the unique identifier for login.
type User struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Mobile    string
	PinCode   string
	State     string
	City      string
	Country   string
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
               email TEXT PRIMARY KEY,
               password TEXT,
               first_name TEXT,
               last_name TEXT,
                mobile TEXT,
                pin_code TEXT,
                state TEXT,
                city TEXT,
                country TEXT
        )`); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) CreateUser(u User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO users(email, password, first_name, last_name, mobile, pin_code, state, city, country)
               VALUES(?,?,?,?,?,?,?,?,?)`,
		u.Email, string(hash), u.FirstName, u.LastName, u.Mobile, u.PinCode, u.State, u.City, u.Country)
	return err
}

func (s *Store) getHash(email string) (string, error) {
	var hash string
	err := s.db.QueryRow(`SELECT password FROM users WHERE email=?`, email).Scan(&hash)
	return hash, err
}

func (s *Store) ValidateUser(email, password string) error {
	hash, err := s.getHash(email)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}

// GetUser returns user details without password.
func (s *Store) GetUser(email string) (*User, error) {
	var u User
	err := s.db.QueryRow(`SELECT email, first_name, last_name, mobile, pin_code, state, city, country FROM users WHERE email=?`, email).
		Scan(&u.Email, &u.FirstName, &u.LastName, &u.Mobile, &u.PinCode, &u.State, &u.City, &u.Country)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) UpdatePassword(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`UPDATE users SET password=? WHERE email=?`, string(hash), email)
	return err
}

// All returns all users without passwords.
func (s *Store) All() ([]User, error) {
	rows, err := s.db.Query(`SELECT email, first_name, last_name, mobile, pin_code, state, city, country FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Email, &u.FirstName, &u.LastName, &u.Mobile, &u.PinCode, &u.State, &u.City, &u.Country); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
