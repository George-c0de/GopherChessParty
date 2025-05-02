package routers

import (
	"net/http"

	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.JSON(http.StatusOK, gin.H{"items": games})
	})
	users.GET("/:game_id", func(c *gin.Context) {
		service := GetService(c)
		_, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		gameID, err := uuid.Parse(c.Param("game_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		games, err := service.GetGameByID(gameID)
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
