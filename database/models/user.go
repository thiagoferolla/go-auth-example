package models

import (
	"database/sql"
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID              uuid.UUID   `json:"id"`
	Name            null.String `json:"name"`
	Email           string      `json:"email"`
	Password        string      `json:"-"`
	Provider        string      `json:"provider"`
	EmailVerifiedAt null.Time   `json:"email_verified_at"`
	Role            string      `json:"role"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type UserRepository interface {
	BeginTransaction() (*sql.Tx, error)
	GetUserByID(id string) (User, error)
	GetUserByEmail(email string) (User, error)
	CreateUser(user *User, transaction *sql.Tx) (*User, error)
	UpdateUser(user *User, transaction *sql.Tx) (*User, error)
	DeleteUser(id string, transaction *sql.Tx) error
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	if err != nil {
		return false
	}

	return true
}

func NewUser(name string, email string, password string, provider string) (*User, error) {
	if len(password) <= 0 {
		return nil, errors.New("password is required")
	} else if !ValidateEmail(email) {
		return nil, errors.New("invalid email")
	}

	newUser := &User{
		ID:       uuid.New(),
		Name:     null.NewString(name, len(name) > 0),
		Email:    email,
		Password: password,
		Provider: provider,
		Role:     "user",
	}

	err := newUser.HashPassword()

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

func (u User) VerifyPassword(pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))

	if err != nil {
		return false
	}

	return true
}
