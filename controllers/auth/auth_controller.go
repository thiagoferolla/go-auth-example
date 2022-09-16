package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/thiagoferolla/go-auth/database/models"
	"github.com/thiagoferolla/go-auth/providers/jwt"
)

type AuthController struct {
	UserRepository models.UserRepository
	JwtProvider    *jwt.JWTProvider
}

func NewAuthController(userRepository models.UserRepository, jwtProvider *jwt.JWTProvider) *AuthController {
	return &AuthController{userRepository, jwtProvider}
}

type CreateUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	IsNewUser bool   `json:"is_new_user"`
	IDToken   string `json:"id_token"`
}

func (controller AuthController) CreateUser(c *gin.Context) {
	payload := CreateUserPayload{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := models.NewUser(payload.Name, payload.Email, payload.Password, "password")

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	_, err = controller.UserRepository.CreateUser(user)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.JwtProvider.GenerateToken(*user)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := CreateUserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		IsNewUser: true,
		IDToken:   token,
	}

	c.JSON(201, response)

	return
}
