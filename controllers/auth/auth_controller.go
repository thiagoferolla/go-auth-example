package auth

import (
	"log"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.NewUser(payload.Name, payload.Email, payload.Password, "password")

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transaction, err := controller.UserRepository.BeginTransaction()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = controller.UserRepository.CreateUser(user, transaction)

	if err != nil {
		transaction.Rollback()
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken := models.NewRefreshToken(user.ID)
	_, err = controller.RefreshTokenRepository.CreateRefreshToken(refreshToken, transaction)

	if err != nil {
		transaction.Rollback()
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.JwtProvider.GenerateToken(*user)

	if err != nil {
		transaction.Rollback()
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = transaction.Commit()

	if err != nil {
		transaction.Rollback()
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := AuthResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		IsNewUser:    true,
		IDToken:      token,
		RefreshToken: refreshToken.Token.String(),
	}

	c.JSON(http.StatusCreated, response)

	return
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (controller AuthController) Login(c *gin.Context) {
	var payload LoginPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(payload.Password) <= 0 || !models.ValidateEmail(payload.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	user, err := controller.UserRepository.GetUserByEmail(payload.Email)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email or password"})
		return
	}

	validPassword := user.VerifyPassword(payload.Password)

	if !validPassword {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	refreshToken := models.NewRefreshToken(user.ID)
	_, err = controller.RefreshTokenRepository.CreateRefreshToken(refreshToken, nil)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.JwtProvider.GenerateToken(user)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := AuthResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		IsNewUser:    false,
		IDToken:      token,
		RefreshToken: refreshToken.Token.String(),
	}

	c.JSON(http.StatusOK, response)

	return
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token"`
}

func (controller AuthController) RefreshToken(c *gin.Context) {
	var payload RefreshTokenPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := controller.RefreshTokenRepository.GetRefreshTokenByToken(payload.RefreshToken)

	if err != nil || len(refreshToken.Owner.String()) <= 0 {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	user, err := controller.UserRepository.GetUserByID(refreshToken.Owner.String())

	if err != nil || len(user.ID.String()) <= 0 {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	token, err := controller.JwtProvider.GenerateToken(user)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := AuthResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		IsNewUser:    false,
		IDToken:      token,
		RefreshToken: refreshToken.Token.String(),
	}

	c.JSON(http.StatusOK, response)

	return
}

func (controller AuthController) Logout(c *gin.Context) {
	var payload RefreshTokenPayload
	user := c.MustGet("user").(models.User)

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := controller.RefreshTokenRepository.GetRefreshTokenByToken(payload.RefreshToken)

	if err != nil || len(refreshToken.Owner.String()) <= 0 || user.ID != refreshToken.Owner {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	err = controller.RefreshTokenRepository.InvalidateToken(refreshToken.Token.String())

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
	c.Abort()

	return
}
