package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/thiagoferolla/go-auth/database"
	"github.com/thiagoferolla/go-auth/routes"
)

func main() {
	databaseConnection, err := database.Connect()

	if err != nil {
		panic(err)
	}

	server := gin.Default()

	routes.NewRouter(server, databaseConnection)

	server.Run()
}