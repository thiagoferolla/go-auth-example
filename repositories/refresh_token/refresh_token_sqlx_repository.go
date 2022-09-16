package refreshtoken

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/database"
	"github.com/thiagoferolla/go-auth/database/models"
)

type RefreshTokenSqlxRepository struct {
	Database *sqlx.DB
}

func NewRefreshTokenSqlxRepository(db *sqlx.DB) *RefreshTokenSqlxRepository {
	return &RefreshTokenSqlxRepository{db}
}

func (r RefreshTokenSqlxRepository) GetRefreshTokenByToken(token string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	err := r.Database.QueryRow("SELECT token, owner, valid, created_at, updated_at FROM refresh_tokens WHERE token = $1", token).
		Scan(&refreshToken.Token, &refreshToken.Owner, &refreshToken.Valid, &refreshToken.CreatedAt, &refreshToken.UpdatedAt)

	return refreshToken, err
}

func (r RefreshTokenSqlxRepository) InvalidateToken(token string) error {
	row := r.Database.QueryRow("UPDATE refresh_tokens SET valid = false WHERE token = $1", token)

	return row.Err()
}

func (r RefreshTokenSqlxRepository) CreateRefreshToken(refreshToken *models.RefreshToken, transaction *sql.Tx) (*models.RefreshToken, error) {
	client := database.ParseClient(r.Database, transaction)

	err := client.QueryRow("INSERT INTO refresh_tokens (token, owner, valid) VALUES ($1, $2, $3) RETURNING token, owner, valid, created_at, updated_at", refreshToken.Token, refreshToken.Owner, refreshToken.Valid).
		Scan(&refreshToken.Token, &refreshToken.Owner, &refreshToken.Valid, &refreshToken.CreatedAt, &refreshToken.UpdatedAt)

	return refreshToken, err
}
