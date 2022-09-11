package user

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/database/models"
)

type UserSqlxRepository struct {
	db *sqlx.DB
}

func NewUserSqlxRepository(db *sqlx.DB) *UserSqlxRepository {
	return &UserSqlxRepository{db}
}

func (r UserSqlxRepository) GetUserByID(id string) (models.User, error) {
	var user models.User

	err := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt,)

	return user, err
}

func (r UserSqlxRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := r.db.QueryRow("SELECT * FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r UserSqlxRepository) CreateUser(user *models.User) (*models.User, error) {
	err := r.db.QueryRow("INSERT INTO users (id, name, email, password, provider) VALUES ($1, $2, $3, $4, $5) RETURNING *", user.ID, user.Name, user.Email, user.Password, user.Provider).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r UserSqlxRepository) UpdateUser(user *models.User) (*models.User, error) {
	err := r.db.QueryRow("UPDATE users SET name = $1, email = $2, password = $3, email_verified_at = $5, updated_at = NOW() WHERE id = $7 RETURNING *", user.Name, user.Email, user.Password, user.EmailVerifiedAt).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r UserSqlxRepository) DeleteUser(id string) error {
	rows, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)

	numberOfRows, _ := rows.RowsAffected()

	if numberOfRows == 0 {
		return errors.New("User not found")
	}

	return err
}