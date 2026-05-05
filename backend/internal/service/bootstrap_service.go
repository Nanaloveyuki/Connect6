package service

import (
	"strings"

	"connect6/backend/internal/config"
	"connect6/backend/internal/constant"
	domainbootstrap "connect6/backend/internal/domain/bootstrap"
)

type BootstrapService struct {
	github config.GitHubConfig
}

func NewBootstrapService(github config.GitHubConfig) *BootstrapService {
	return &BootstrapService{github: github}
}

func (s *BootstrapService) GetConfig() domainbootstrap.Config {
	capability := domainbootstrap.GitHubCapability{
		Mode:        "limited",
		Description: "GitHub token missing. Live requests are available with stricter limits.",
	}

	if strings.TrimSpace(s.github.Token) != "" {
		capability = domainbootstrap.GitHubCapability{
			Mode:        "authenticated",
			Description: "GitHub token configured. Live requests use authenticated limits.",
		}
	}

	return domainbootstrap.Config{
		Search: domainbootstrap.SearchConfig{
			DefaultSource:     constant.DefaultBootstrapSource,
			DefaultTarget:     constant.DefaultBootstrapTarget,
			SourcePlaceholder: constant.DefaultSourcePlaceholder,
			TargetPlaceholder: constant.DefaultTargetPlaceholder,
			DefaultMaxDepth:   constant.DefaultBootstrapMaxDepth,
		},
		Controls: domainbootstrap.ControlsConfig{
			MaxDepthMin:       constant.DefaultBootstrapMaxDepthMin,
			MaxDepthMax:       constant.DefaultBootstrapMaxDepthMax,
			DefaultRepulsion:  constant.DefaultRepulsion,
			RepulsionMin:      constant.DefaultRepulsionMin,
			RepulsionMax:      constant.DefaultRepulsionMax,
			DefaultAttraction: constant.DefaultAttraction,
			AttractionMin:     constant.DefaultAttractionMin,
			AttractionMax:     constant.DefaultAttractionMax,
			DefaultNodeSize:   constant.DefaultNodeSize,
			NodeSizeMin:       constant.DefaultNodeSizeMin,
			NodeSizeMax:       constant.DefaultNodeSizeMax,
		},
		Graph: domainbootstrap.GraphConfig{
			InitialZoom:           constant.DefaultInitialZoom,
			MinZoomRatio:          constant.DefaultMinZoomRatio,
			MaxZoomRatio:          constant.DefaultMaxZoomRatio,
			LayoutIterations:      constant.DefaultForceAtlasIterations,
			PlaybackInterval:      constant.DefaultPlaybackIntervalMS,
			StatusCardMaxWidthRem: constant.DefaultGraphStatusCardMaxWidth,
			GitHubCapability:      capability,
			RelationLegend: []domainbootstrap.RelationLegendItem{
				{Key: constant.RelationFollows, Label: "Follows", Color: "#8b7cf6"},
				{Key: constant.RelationMaintainerBridge, Label: "Maintainer bridge", Color: "#2f7df6"},
				{Key: constant.RelationFrameworkCircle, Label: "Framework circle", Color: "#18a999"},
				{Key: constant.RelationPackageLine, Label: "Package line", Color: "#f39c56"},
			},
			EnablePan:             constant.DefaultEnablePan,
			EnableZoom:            constant.DefaultEnableZoom,
			EnableAnimations:      constant.DefaultEnableAnimations,
			EnableSidebar:         constant.DefaultEnableSidebar,
		},
	}
}
