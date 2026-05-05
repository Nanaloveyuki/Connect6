package http

import (
	stdhttp "net/http"

	"connect6/backend/internal/dto"
	"connect6/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type BootstrapHandler struct {
	bootstrapService *service.BootstrapService
}

func NewBootstrapHandler(bootstrapService *service.BootstrapService) *BootstrapHandler {
	return &BootstrapHandler{bootstrapService: bootstrapService}
}

func (h *BootstrapHandler) Get(c *gin.Context) {
	config := h.bootstrapService.GetConfig()
	legend := make([]dto.RelationLegendItem, 0, len(config.Graph.RelationLegend))
	for _, item := range config.Graph.RelationLegend {
		legend = append(legend, dto.RelationLegendItem{
			Key:   item.Key,
			Label: item.Label,
			Color: item.Color,
		})
	}

	c.JSON(stdhttp.StatusOK, dto.BootstrapResponse{
		Search: dto.SearchBootstrap{
			DefaultSource:     config.Search.DefaultSource,
			DefaultTarget:     config.Search.DefaultTarget,
			SourcePlaceholder: config.Search.SourcePlaceholder,
			TargetPlaceholder: config.Search.TargetPlaceholder,
			DefaultMaxDepth:   config.Search.DefaultMaxDepth,
		},
		Controls: dto.ControlsBootstrap{
			MaxDepthMin:       config.Controls.MaxDepthMin,
			MaxDepthMax:       config.Controls.MaxDepthMax,
			DefaultRepulsion:  config.Controls.DefaultRepulsion,
			RepulsionMin:      config.Controls.RepulsionMin,
			RepulsionMax:      config.Controls.RepulsionMax,
			DefaultAttraction: config.Controls.DefaultAttraction,
			AttractionMin:     config.Controls.AttractionMin,
			AttractionMax:     config.Controls.AttractionMax,
			DefaultNodeSize:   config.Controls.DefaultNodeSize,
			NodeSizeMin:       config.Controls.NodeSizeMin,
			NodeSizeMax:       config.Controls.NodeSizeMax,
		},
		Graph: dto.GraphBootstrap{
			InitialZoom:          config.Graph.InitialZoom,
			MinZoomRatio:         config.Graph.MinZoomRatio,
			MaxZoomRatio:         config.Graph.MaxZoomRatio,
			LayoutIterations:     config.Graph.LayoutIterations,
			PlaybackInterval:     config.Graph.PlaybackInterval,
			StatusCardMaxWidthRem: config.Graph.StatusCardMaxWidthRem,
			GitHubCapability: dto.GitHubCapability{
				Mode:        config.Graph.GitHubCapability.Mode,
				Description: config.Graph.GitHubCapability.Description,
			},
			RelationLegend:       legend,
			EnablePan:            config.Graph.EnablePan,
			EnableZoom:           config.Graph.EnableZoom,
			EnableAnimations:     config.Graph.EnableAnimations,
			EnableSidebar:        config.Graph.EnableSidebar,
		},
	})
}
