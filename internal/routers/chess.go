package routers

import (
	"GopherChessParty/internal/middleware"
	"GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
)

func addChessRoute(rg *gin.RouterGroup, service services.IService) {
	users := rg.Group("/chess")
	users.Use(middleware.JWTAuthMiddleware(service))
	users.GET("/", func(c *gin.Context) {
		// TODO GET ALL GAMES
	})

	users.POST("/", func(c *gin.Context) {
		// TODO CREATE GAME
	})
}
