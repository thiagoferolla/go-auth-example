package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/thiagoferolla/go-auth/database"
	"github.com/thiagoferolla/go-auth/providers/secret"
	"github.com/thiagoferolla/go-auth/routes"
)

func main() {
	secretProvider := secret.NewMockSecretProvider()

	databaseConnection, err := database.Connect(secretProvider)

	if err != nil {
		panic(err)
	}

	server := gin.Default()

	routes.NewRouter(server, databaseConnection)

	server.Run()
}
