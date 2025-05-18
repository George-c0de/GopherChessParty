package routers

import (
	"net/http"

	"GopherChessParty/internal/dto"
	"github.com/gin-gonic/gin"
)

func addAuthRoutes(rg *gin.RouterGroup) {
	// Публичный маршрут для логина
	rg.POST("/login", func(c *gin.Context) {
		data, err := BindJSON[dto.AuthenticateUser](c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		service := GetService(c)

		userId, ok := service.ValidPassword(*data)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Valid Password"})
			return
		}

		tokens, err := service.GenerateTokenPair(*userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		})
	})

	// Публичный маршрут для регистрации
	rg.POST("/register", func(c *gin.Context) {
		data, err := BindJSON[dto.CreateUser](c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		service := GetService(c)

		user, err := service.RegisterUser(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tokens, err := service.GenerateTokenPair(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		})
	})

	// Публичный маршрут для обновления токена
	rg.POST("/refresh", func(c *gin.Context) {
		var data struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		service := GetService(c)
		tokens, err := service.RefreshAccessToken(data.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		})
	})
}
