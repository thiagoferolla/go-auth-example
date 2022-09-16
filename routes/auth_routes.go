package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/controllers/auth"
	"github.com/thiagoferolla/go-auth/providers/jwt"
	"github.com/thiagoferolla/go-auth/repositories/user"
)

func RegisterAuthRoutes(server *gin.Engine, database *sqlx.DB, jwtProvider *jwt.JWTProvider) {
	group := server.Group("/auth/v1")

	authController := auth.NewAuthController(
		user.NewUserSqlxRepository(database),
		jwtProvider,
	)

	group.POST("/sign-in", authController.CreateUser)
}
