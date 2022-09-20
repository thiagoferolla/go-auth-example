package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/providers/secret"
)

var host = os.Getenv("DB_HOST")
var port = os.Getenv("DB_PORT")
var user = os.Getenv("DB_USER")
var password = os.Getenv("DB_PASSWORD")
var name = os.Getenv("DB_NAME")

func GetConnectionString(secretProvider secret.SecretProvider) string {
	host, err := secretProvider.Get("DB_HOST")

	if err != nil {
		panic(err)
	}

	port, err := secretProvider.Get("DB_PORT")

	if err != nil {
		port = "5432"
	}

	user, err := secretProvider.Get("DB_USER")

	if err != nil {
		panic(err)
	}

	password, err := secretProvider.Get("DB_PASSWORD")

	if err != nil {
		panic(err)
	}

	name, err := secretProvider.Get("DB_NAME")

	if err != nil {
		panic(err)
	}


	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
}

func Connect(secretProvider secret.SecretProvider) (*sqlx.DB, error) {
	var err error

	connectionString := GetConnectionString(secretProvider)

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
