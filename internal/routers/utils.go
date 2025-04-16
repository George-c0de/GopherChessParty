package routers

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"net/http"

	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
)

func GetService(c *gin.Context) interfaces.IService {
	svc, exists := c.Get("service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service not found"})
		return nil
	}
	// Приводим к нужному типу
	return svc.(*services.Service)
}
func GetUserID(c *gin.Context) *uuid.UUID {
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return nil
	}
	return userId.(*uuid.UUID)
}
func BindJSON[T any](c *gin.Context) (*T, error) {
	var data T
	if err := c.BindJSON(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
