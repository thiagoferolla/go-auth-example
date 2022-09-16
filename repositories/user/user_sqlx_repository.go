package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/database"
	"github.com/thiagoferolla/go-auth/database/models"
)

type UserSqlxRepository struct {
	Database *sqlx.DB
}

func NewUserSqlxRepository(db *sqlx.DB) *UserSqlxRepository {
	return &UserSqlxRepository{db}
}

func (r UserSqlxRepository) BeginTransaction() (*sql.Tx, error) {
	c := context.Background()

	return r.Database.BeginTx(c, nil)
}

func (r UserSqlxRepository) GetUserByID(id string) (models.User, error) {
	var user models.User

	err := r.Database.QueryRow("SELECT id, name, email, password, provider, email_verified_at, role, created_at, updated_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.EmailVerifiedAt, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r UserSqlxRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := r.Database.QueryRow("SELECT id, name, email, password, provider, email_verified_at, role, created_at, updated_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.EmailVerifiedAt, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r UserSqlxRepository) CreateUser(user *models.User, transaction *sql.Tx) (*models.User, error) {
	client := database.ParseClient(r.Database, transaction)

	err := client.QueryRow("INSERT INTO users (id, name, email, password, provider, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, email, password, provider, email_verified_at, role, created_at, updated_at", user.ID, user.Name, user.Email, user.Password, user.Provider, user.Role).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.EmailVerifiedAt, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r UserSqlxRepository) UpdateUser(user *models.User, transaction *sql.Tx) (*models.User, error) {
	client := database.ParseClient(r.Database, transaction)

	err := client.QueryRow("UPDATE users SET name = $1, email = $2, password = $3, email_verified_at = $5, role = $6, updated_at = NOW() WHERE id = $7 RETURNING id, name, email, password, provider, email_verified_at, role, created_at, updated_at", user.Name, user.Email, user.Password, user.EmailVerifiedAt, user.Role).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.EmailVerifiedAt, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r UserSqlxRepository) DeleteUser(id string, transaction *sql.Tx) error {
	client := database.ParseClient(r.Database, transaction)

	rows, err := client.Exec("DELETE FROM users WHERE id = $1", id)

	numberOfRows, _ := rows.RowsAffected()

	if numberOfRows == 0 {
		return errors.New("User not found")
	}

	return err
}
