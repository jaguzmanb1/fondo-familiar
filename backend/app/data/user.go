package data

import (
	"database/sql"
	"fmt"

	"github.com/hashicorp/go-hclog"
)

// ErrUserNotFound is raised when a user is not found
var ErrUserNotFound = fmt.Errorf("User not found")

// ErrCreditNotFound is raised when a user is not found
var ErrCreditNotFound = fmt.Errorf("Credit not found")

// User describes a user
type User struct {
	ID    int
	Rol   int
	Email string
}

// UserService does
type UserService struct {
	DB *sql.DB
	l  hclog.Logger
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB, l hclog.Logger) *UserService {
	return &UserService{db, l}
}

// UserExists return true if an specific user id exists
func (u *UserService) UserExists(id int) (bool, error) {
	rows, err := u.DB.Query("SELECT * from usuario where id = ?", id)
	if err != nil {
		return false, ErrUserNotFound
	}

	for rows.Next() {
		return true, nil
	}

	return false, ErrUserNotFound
}

// CreditExists return true if an specific credit id exists
func (u *UserService) CreditExists(id int) (bool, error) {
	rows, err := u.DB.Query("SELECT * from creditos where id = ?", id)
	if err != nil {
		return false, ErrCreditNotFound
	}

	for rows.Next() {
		return true, nil
	}

	return false, ErrCreditNotFound
}
