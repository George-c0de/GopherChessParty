package routers

import (
	"net/http"

	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/middleware"
	"GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup, service services.IService) {
	users := rg.Group("/users")
	users.Use(middleware.JWTAuthMiddleware(service))
	users.GET("/", func(c *gin.Context) {
		service := GetService(c)
		users, err := service.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": users})
		return
	})

	users.POST("/", func(c *gin.Context) {
		data, err := BindJSON[dto.CreateUser](c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		service := GetService(c)
		user, err := service.CreateUser(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not create user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"item": user})
		return
	})
}
