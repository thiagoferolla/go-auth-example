package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var host = os.Getenv("DB_HOST")
var port = os.Getenv("DB_PORT")
var user = os.Getenv("DB_USER")
var password = os.Getenv("DB_PASSWORD")
var name = os.Getenv("DB_NAME")

func GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
}

func Connect() (*sqlx.DB, error) {
	var err error

	connectionString := GetConnectionString()

	connection, err := sqlx.Connect("pgx", connectionString)

	return connection, err
}

type QueryClient interface {
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

func ParseClient(database *sqlx.DB, transaction *sql.Tx) QueryClient {
	if transaction != nil {
		return transaction
	}

	return database
}
