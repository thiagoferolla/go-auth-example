package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/controllers/auth"
	"github.com/thiagoferolla/go-auth/providers/jwt"
	refreshtoken "github.com/thiagoferolla/go-auth/repositories/refresh_token"
	"github.com/thiagoferolla/go-auth/repositories/user"
)

func RegisterAuthRoutes(server *gin.Engine, database *sqlx.DB, jwtProvider *jwt.JWTProvider) {
	group := server.Group("/auth/v1")

	authController := auth.NewAuthController(
		user.NewUserSqlxRepository(database),
		refreshtoken.NewRefreshTokenSqlxRepository(database),
		jwtProvider,
	)

	group.POST("/sign-in", authController.CreateUser)
}
