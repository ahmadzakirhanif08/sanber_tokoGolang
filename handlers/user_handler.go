package handlers

import (
	"net/http"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/configs"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/models"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/middlewares"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var input models.RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed encrypting password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// 2. Simpan ke Database
	result := configs.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "sign up success", "username": user.Username})
}

func LoginHandler(c *gin.Context) {
	var input models.LoginRequest
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := configs.DB.Where("username = ?", input.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password is not correct"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password is not correct"})
		return
	}

	// 3. Generate JWT Token
	token, err := middlewares.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login success", "token": token})
}