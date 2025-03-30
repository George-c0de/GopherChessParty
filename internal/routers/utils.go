package routers

import (
	"GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
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

func BindJSON[T any](c *gin.Context) (*T, error) {
	var data T
	if err := c.BindJSON(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
