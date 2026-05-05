package bootstrap

type Config struct {
	Search   SearchConfig
	Controls ControlsConfig
	Graph    GraphConfig
}

type SearchConfig struct {
	DefaultSource     string
	DefaultTarget     string
	SourcePlaceholder string
	TargetPlaceholder string
	DefaultMaxDepth   int
}

type ControlsConfig struct {
	MaxDepthMin       int
	MaxDepthMax       int
	DefaultRepulsion  int
	RepulsionMin      int
	RepulsionMax      int
	DefaultAttraction int
	AttractionMin     int
	AttractionMax     int
	DefaultNodeSize   int
	NodeSizeMin       int
	NodeSizeMax       int
}

type GraphConfig struct {
	InitialZoom           float64
	MinZoomRatio          float64
	MaxZoomRatio          float64
	LayoutIterations      int
	PlaybackInterval      int
	StatusCardMaxWidthRem int
	GitHubCapability      GitHubCapability
	RelationLegend        []RelationLegendItem
	EnablePan             bool
	EnableZoom            bool
	EnableAnimations      bool
	EnableSidebar         bool
}

type GitHubCapability struct {
	Mode        string
	Description string
}

type RelationLegendItem struct {
	Key   string
	Label string
	Color string
}
