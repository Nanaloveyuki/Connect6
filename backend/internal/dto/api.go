package dto

import "connect6/backend/internal/domain/graph"

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Login      string `json:"login"`
	Name       string `json:"name"`
	AvatarURL  string `json:"avatarUrl,omitempty"`
	ProfileURL string `json:"profileUrl,omitempty"`
	Bio        string `json:"bio,omitempty"`
}

type ShortestPathRequest struct {
	Source   string `json:"source" binding:"required"`
	Target   string `json:"target" binding:"required"`
	MaxDepth int    `json:"maxDepth"`
}

type BootstrapResponse struct {
	Search   SearchBootstrap   `json:"search"`
	Controls ControlsBootstrap `json:"controls"`
	Graph    GraphBootstrap    `json:"graph"`
}

type SearchBootstrap struct {
	DefaultSource      string `json:"defaultSource"`
	DefaultTarget      string `json:"defaultTarget"`
	SourcePlaceholder  string `json:"sourcePlaceholder"`
	TargetPlaceholder  string `json:"targetPlaceholder"`
	DefaultMaxDepth    int    `json:"defaultMaxDepth"`
}

type ControlsBootstrap struct {
	MaxDepthMin          int `json:"maxDepthMin"`
	MaxDepthMax          int `json:"maxDepthMax"`
	DefaultRepulsion     int `json:"defaultRepulsion"`
	RepulsionMin         int `json:"repulsionMin"`
	RepulsionMax         int `json:"repulsionMax"`
	DefaultAttraction    int `json:"defaultAttraction"`
	AttractionMin        int `json:"attractionMin"`
	AttractionMax        int `json:"attractionMax"`
	DefaultNodeSize      int `json:"defaultNodeSize"`
	NodeSizeMin          int `json:"nodeSizeMin"`
	NodeSizeMax          int `json:"nodeSizeMax"`
}

type GraphBootstrap struct {
	InitialZoom           float64              `json:"initialZoom"`
	MinZoomRatio          float64              `json:"minZoomRatio"`
	MaxZoomRatio          float64              `json:"maxZoomRatio"`
	LayoutIterations      int                  `json:"layoutIterations"`
	PlaybackInterval      int                  `json:"playbackInterval"`
	StatusCardMaxWidthRem int                  `json:"statusCardMaxWidthRem"`
	GitHubCapability      GitHubCapability     `json:"githubCapability"`
	RelationLegend        []RelationLegendItem `json:"relationLegend"`
	EnablePan             bool                 `json:"enablePan"`
	EnableZoom            bool                 `json:"enableZoom"`
	EnableAnimations      bool                 `json:"enableAnimations"`
	EnableSidebar         bool                 `json:"enableSidebar"`
}

type GitHubCapability struct {
	Mode        string `json:"mode"`
	Description string `json:"description"`
}

type RelationLegendItem struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Color string `json:"color"`
}

type PathResponse = graph.Path
