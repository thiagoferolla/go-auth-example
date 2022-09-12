package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/controllers/auth"
	"github.com/thiagoferolla/go-auth/repositories/user"
)

func RegisterAuthRoutes(server *gin.Engine, database *sqlx.DB) {
	group := server.Group("/auth")

	authController := auth.NewAuthController(
		user.NewUserSqlxRepository(database),
	)

	group.POST("/sign-in", authController.CreateUser)
}