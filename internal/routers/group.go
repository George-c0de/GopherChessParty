package routers

import (
	"github.com/George-c0de/GopherChessParty/internal/middleware"
	"github.com/George-c0de/GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
)

func GetRoutes(service *services.Service) *gin.Engine {
	router := gin.Default()
	// Применение middleware для добавления сервиса в контекст
	router.Use(middleware.ServiceMiddleware(service))
	v1 := router.Group("/v1")
	addUserRoutes(v1)
	return router
}
