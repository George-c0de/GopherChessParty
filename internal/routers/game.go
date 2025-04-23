package routers

import (
	"net/http"

	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"github.com/gin-gonic/gin"
)

func addChessRoute(rg *gin.RouterGroup, service interfaces.IService) {
	users := rg.Group("/chess")
	users.Use(middleware.JWTAuthMiddleware(service))
	users.GET("/", func(c *gin.Context) {
		service := GetService(c)
		userId, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		games, err := service.GetGamesByUserID(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, games)
	})

	users.POST("/", func(c *gin.Context) {
		// TODO CREATE GAME
	})
}
