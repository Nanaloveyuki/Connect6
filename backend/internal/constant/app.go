package constant

const (
	DefaultEnvironment             = "development"
	DefaultHTTPPort                = "8080"
	DefaultGitHubBaseURL           = "https://api.github.com"
	DefaultHTTPReadTimeoutSeconds  = 10
	DefaultHTTPWriteTimeoutSeconds = 10
	DefaultGitHubTimeoutSeconds    = 15
	DefaultGitHubFollowingLimit    = 60
	DefaultGitHubOrgLimit          = 20
	DefaultGraphMaxDepth           = 6
	DefaultBootstrapSource         = "sveltejs"
	DefaultBootstrapTarget         = "gin-gonic"
	DefaultSourcePlaceholder       = "github-user-a"
	DefaultTargetPlaceholder       = "github-user-b"
	DefaultBootstrapMaxDepth       = 4
	DefaultBootstrapMaxDepthMin    = 1
	DefaultBootstrapMaxDepthMax    = 6
	DefaultRepulsion               = 36
	DefaultRepulsionMin            = 10
	DefaultRepulsionMax            = 80
	DefaultAttraction              = 14
	DefaultAttractionMin           = 5
	DefaultAttractionMax           = 30
	DefaultNodeSize                = 14
	DefaultNodeSizeMin             = 8
	DefaultNodeSizeMax             = 24
	DefaultInitialZoom             = 1.0
	DefaultMinZoomRatio            = 0.2
	DefaultMaxZoomRatio            = 4.0
	DefaultForceAtlasIterations    = 120
	DefaultPlaybackIntervalMS      = 1200
	DefaultGraphStatusCardMaxWidth = 32
	DefaultEnablePan               = true
	DefaultEnableZoom              = true
	DefaultEnableAnimations        = true
	DefaultEnableSidebar           = true
	RelationMaintainerBridge       = "maintainer-bridge"
	RelationFrameworkCircle        = "framework-circle"
	RelationPackageLine            = "package-line"
	RelationFollows                = "follows"
)

var DefaultAllowedOrigins = []string{
	"http://localhost:5173",
}
