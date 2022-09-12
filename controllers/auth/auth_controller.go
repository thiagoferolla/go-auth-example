package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/thiagoferolla/go-auth/database/models"
)

type AuthController struct {
	UserRepository models.UserRepository
}

func NewAuthController(userRepository models.UserRepository) *AuthController {
	return &AuthController{userRepository}
}

type CreateUserPayload struct {
	Name string
	Email string
	Password string
	Provider string
}

func (controller AuthController) CreateUser(c *gin.Context) {
	payload := CreateUserPayload{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := models.NewUser(payload.Name, payload.Email, payload.Password, payload.Provider)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return 
	}

	_, err = controller.UserRepository.CreateUser(user)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)

	return
}