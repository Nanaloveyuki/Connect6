package router

import (
	handlerhttp "connect6/backend/internal/handler/http"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine, healthHandler *handlerhttp.HealthHandler, bootstrapHandler *handlerhttp.BootstrapHandler, userHandler *handlerhttp.UserHandler, graphHandler *handlerhttp.GraphHandler) {
	api := engine.Group("/api/v1")
	api.GET("/health", healthHandler.Get)
	api.GET("/bootstrap", bootstrapHandler.Get)
	api.GET("/users/:username", userHandler.GetByUsername)
	api.POST("/graph/path", graphHandler.FindShortestPath)
}
