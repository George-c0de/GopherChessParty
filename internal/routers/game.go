package routers

import (
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"github.com/gin-gonic/gin"
)

func addChessRoute(rg *gin.RouterGroup, service interfaces.IService) {
	users := rg.Group("/chess")
	users.Use(middleware.JWTAuthMiddleware(service))
	users.GET("/", func(c *gin.Context) {
		service := GetService(c)
		userId := *GetUserID(c)
		service.GetGames(userId)

	})

	users.POST("/", func(c *gin.Context) {
		// TODO CREATE GAME
	})
}
