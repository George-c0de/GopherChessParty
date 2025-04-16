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

		token, err := service.GenerateToken(*userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}
