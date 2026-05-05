package http

import (
	stdhttp "net/http"

	"connect6/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Get(c *gin.Context) {
	c.JSON(stdhttp.StatusOK, dto.HealthResponse{
		Status:  "ok",
		Service: "connect6-backend",
	})
}

