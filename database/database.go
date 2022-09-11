package database

import (
	"fmt"
	"os"

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
	connectionString := GetConnectionString()

	return sqlx.Connect("postgres", connectionString)
}