package routers

import (
	"GopherChessParty/internal/middleware"
	"GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
)

func GetRoutes(service services.IService) *gin.Engine {
	router := gin.Default()
	// Применение middleware для добавления сервиса в контекст
	router.Use(middleware.ServiceMiddleware(service))
	v1 := router.Group("/v1")
	addAuthRoutes(v1)
	addUserRoutes(v1, service)
	addChessRoute(v1, service)
	return router
}
