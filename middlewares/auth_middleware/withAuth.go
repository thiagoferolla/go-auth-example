package auth_middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thiagoferolla/go-auth/database/models"
	"github.com/thiagoferolla/go-auth/providers/jwt"
)

type WithAuthMiddleware struct {
	UserRepository models.UserRepository
	JwtProvider    jwt.JWTProvider
}

func NewWithAuthMiddleware(userRepository models.UserRepository, jwtProvider jwt.JWTProvider) *WithAuthMiddleware {
	return &WithAuthMiddleware{userRepository, jwtProvider}
}

func (auth WithAuthMiddleware) WithAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) <= 0 {
			c.AbortWithStatusJSON(403, gin.H{"error": "Not authorized"})
			return
		}

		tokenMap := strings.Split(authHeader, "Bearer")

		if len(tokenMap) < 2 {
			c.AbortWithStatusJSON(403, gin.H{"error": "Not authorized"})
			return
		}

		bearerToken := strings.TrimSpace(tokenMap[1])

		claims, err := auth.JwtProvider.ValidateToken(bearerToken)

		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "Not authorized"})
			return
		}

		user, err := auth.UserRepository.GetUserByID(claims.ID.String())

		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(403, gin.H{"error": "Not authorized"})
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
