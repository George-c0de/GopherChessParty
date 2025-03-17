package routers

import (
	"github.com/George-c0de/GopherChessParty/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetService(c *gin.Context) *services.Service {
	svc, exists := c.Get("service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service not found"})
		return nil
	}
	// Приводим к нужному типу
	return svc.(*services.Service)
}
func addUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")

	users.GET("/", func(c *gin.Context) {
		service := GetService(c)
		if service == nil {
			return
		}
		users := service.GetUsers()
		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	users.POST("/", func(c *gin.Context) {
		service := GetService(c)
		if service == nil {
			return
		}
	})
}
