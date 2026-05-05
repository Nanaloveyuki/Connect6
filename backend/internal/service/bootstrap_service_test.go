package service

import (
	"testing"

	"connect6/backend/internal/config"
)

func TestBootstrapServiceReturnsFrontendConfig(t *testing.T) {
	t.Parallel()

	service := NewBootstrapService(config.GitHubConfig{})
	config := service.GetConfig()

	if config.Search.DefaultSource == "" || config.Search.DefaultTarget == "" {
		t.Fatal("expected bootstrap search defaults to be populated")
	}

	if config.Controls.RepulsionMin >= config.Controls.RepulsionMax {
		t.Fatal("expected valid repulsion range")
	}

	if !config.Graph.EnableAnimations || !config.Graph.EnableZoom {
		t.Fatal("expected graph interaction defaults to be enabled")
	}

	if config.Graph.LayoutIterations <= 0 || config.Graph.PlaybackInterval <= 0 {
		t.Fatal("expected graph runtime values to be populated")
	}

	if config.Graph.MinZoomRatio <= 0 || config.Graph.MaxZoomRatio <= config.Graph.MinZoomRatio {
		t.Fatal("expected valid zoom ratio range")
	}

	if len(config.Graph.RelationLegend) == 0 {
		t.Fatal("expected relation legend to be populated")
	}

	if config.Graph.GitHubCapability.Mode != "limited" {
		t.Fatal("expected limited capability mode without token")
	}
}

func TestBootstrapServiceMarksAuthenticatedGitHubCapability(t *testing.T) {
	t.Parallel()

	service := NewBootstrapService(config.GitHubConfig{Token: "token-present"})
	config := service.GetConfig()

	if config.Graph.GitHubCapability.Mode != "authenticated" {
		t.Fatalf("expected authenticated capability mode, got %s", config.Graph.GitHubCapability.Mode)
	}
}
