package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Provider string `json:"provider"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	GetUserByID(id uuid.UUID) (User, error)
	GetUserByEmail(email string) (User, error)
	CreateUser(user *User) (User, error)
	UpdateUser(user *User) (User, error)
	DeleteUser(id uuid.UUID) error
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