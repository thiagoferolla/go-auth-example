package auth

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/thiagoferolla/go-auth/database/models"
	"github.com/thiagoferolla/go-auth/providers/jwt"
)

type AuthController struct {
	UserRepository         models.UserRepository
	RefreshTokenRepository models.RefreshTokenRepository
	JwtProvider            *jwt.JWTProvider
}

func NewAuthController(userRepository models.UserRepository, refreshTokenRepository models.RefreshTokenRepository, jwtProvider *jwt.JWTProvider) *AuthController {
	return &AuthController{userRepository, refreshTokenRepository, jwtProvider}
}

type CreateUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	IsNewUser    bool   `json:"is_new_user"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

func (controller AuthController) CreateUser(c *gin.Context) {
	payload := CreateUserPayload{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := models.NewUser(payload.Name, payload.Email, payload.Password, "password")

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transaction, err := controller.UserRepository.BeginTransaction()

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	_, err = controller.UserRepository.CreateUser(user, transaction)

	if err != nil {
		transaction.Rollback()
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	refreshToken := models.NewRefreshToken(user.ID)
	_, err = controller.RefreshTokenRepository.CreateRefreshToken(refreshToken, transaction)

	if err != nil {
		transaction.Rollback()
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.JwtProvider.GenerateToken(*user)

	if err != nil {
		transaction.Rollback()
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = transaction.Commit()

	if err != nil {
		transaction.Rollback()
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := AuthResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		IsNewUser:    true,
		IDToken:      token,
		RefreshToken: refreshToken.Token.String(),
	}

	c.JSON(201, response)

	return
}

type LoginPayload struct {
	Email string 	`json:"email"`
	Password string 	`json:"password"`
}

func (controller AuthController) Login(c *gin.Context) {
	var payload LoginPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(payload.Password) <= 0 || !models.ValidateEmail(payload.Email) {
		c.JSON(400, gin.H{"error": "Invalid email or password"})
		return
	}

	user, err := controller.UserRepository.GetUserByEmail(payload.Email)

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Invalid email or password"})
		return
	}

	validPassword := user.VerifyPassword(payload.Password)

	if !validPassword {
		log.Println(err)
		c.JSON(400, gin.H{"error": "Invalid email or password"})
		return
	}

	refreshToken := models.NewRefreshToken(user.ID)
	_, err = controller.RefreshTokenRepository.CreateRefreshToken(refreshToken, nil)

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.JwtProvider.GenerateToken(user)

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := AuthResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		IsNewUser:    true,
		IDToken:      token,
		RefreshToken: refreshToken.Token.String(),
	}

	c.JSON(200, response)

	return
}

