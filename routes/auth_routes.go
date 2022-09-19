package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/controllers/auth"
	"github.com/thiagoferolla/go-auth/middlewares/auth_middleware"
	"github.com/thiagoferolla/go-auth/providers/cache"
	"github.com/thiagoferolla/go-auth/providers/email"
	"github.com/thiagoferolla/go-auth/providers/jwt"
	refreshtoken "github.com/thiagoferolla/go-auth/repositories/refresh_token"
	"github.com/thiagoferolla/go-auth/repositories/user"
)

func RegisterAuthRoutes(server *gin.Engine, database *sqlx.DB, jwtProvider jwt.JWTProvider, emailProvider email.EmailProvider, cacheProvider cache.CacheProvider) {
	group := server.Group("/auth/v1")

	authController := auth.NewAuthController(
		user.NewUserSqlxRepository(database),
		refreshtoken.NewRefreshTokenSqlxRepository(database),
		jwtProvider,
		emailProvider,
		cacheProvider,
	)

	group.POST("/sign_in", authController.CreateUser)
	group.POST("/login", authController.Login)
	group.POST("/refresh_token", authController.RefreshToken)
	group.POST("/reset_password", authController.SendPasswordReset)

	authMiddleware := auth_middleware.NewWithAuthMiddleware(user.NewUserSqlxRepository(database), jwtProvider)

	withAuthRoutes := group.Group("/")
	withAuthRoutes.Use(authMiddleware.WithAuth())
	withAuthRoutes.POST("/logout", authController.Logout)
}
