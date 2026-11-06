package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/thiagoferolla/go-auth/database"
	"github.com/thiagoferolla/go-auth/providers/jwt"
	"github.com/thiagoferolla/go-auth/providers/secret"
	"github.com/thiagoferolla/go-auth/routes"
)

func main() {
	secretProvider := secret.NewMockSecretProvider()

	jwtSecret, err := secretProvider.Get("JWT_SECRET")

	if err != nil {
		panic(err)
	}

	jwtProvider := jwt.NewBaseProvider([]byte(jwtSecret))

	databaseConnection, err := database.Connect(secretProvider)

	if err != nil {
		panic(err)
	}

	server := gin.Default()

	routes.NewRouter(server, databaseConnection, jwtProvider)

	server.Run()
}
