package http

import (
	"errors"
	stdhttp "net/http"

	"connect6/backend/internal/domain/graph"
	"connect6/backend/internal/dto"
	"connect6/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type GraphHandler struct {
	pathService *service.PathService
}

func NewGraphHandler(pathService *service.PathService) *GraphHandler {
	return &GraphHandler{pathService: pathService}
}

func (h *GraphHandler) FindShortestPath(c *gin.Context) {
	var request dto.ShortestPathRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(stdhttp.StatusBadRequest, dto.ErrorResponse{Message: "invalid path request payload"})
		return
	}

	path, err := h.pathService.FindShortestPath(c.Request.Context(), request.Source, request.Target, request.MaxDepth)
	if err != nil {
		statusCode := stdhttp.StatusInternalServerError
		if errors.Is(err, graph.ErrPathNotFound) {
			statusCode = stdhttp.StatusNotFound
		}

		c.JSON(statusCode, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(stdhttp.StatusOK, dto.PathResponse(path))
}

