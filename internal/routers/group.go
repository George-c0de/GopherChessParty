package routers

import (
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"github.com/gin-gonic/gin"
)

func GetRoutes(service interfaces.IService, log interfaces.ILogger) *gin.Engine {
	router := gin.Default()
	// Применение middleware для добавления сервиса в контекст
	router.Use(middleware.ServiceMiddleware(service))
	v1 := router.Group("/v1")
	AddWebSocket(v1, service, log)
	addAuthRoutes(v1)
	addUserRoutes(v1, service)
	addChessRoute(v1, service)
	return router
}
