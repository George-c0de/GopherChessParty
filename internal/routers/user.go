package routers

import (
	"net/http"

	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup, service interfaces.IService) {
	users := rg.Group("/users")
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
		}

		c.JSON(http.StatusOK, gin.H{"item": user})
	})
	users.Use(middleware.JWTAuthMiddleware(service))
	users.GET("/", func(c *gin.Context) {
		service := GetService(c)
		users, err := service.Users()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": users})
	})
	users.GET("/me", func(c *gin.Context) {
		service := GetService(c)
		userId, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		user, err := service.UserByID(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"item": user})
	})
}
