package bootstrap

import (
	"net/http"
	"strings"
	"time"

	"connect6/backend/internal/config"
	handlerhttp "connect6/backend/internal/handler/http"
	"connect6/backend/internal/repository/memory"
	"connect6/backend/internal/router"
	"connect6/backend/internal/service"
	"connect6/backend/internal/transport/githubapi"

	"github.com/gin-gonic/gin"
)

type Application struct {
	config config.Config
	engine *gin.Engine
}

func NewApplication(configDir string) (*Application, error) {
	cfg, err := config.Load(configDir)
	if err != nil {
		return nil, err
	}

	gin.SetMode(resolveGinMode(cfg.App.Env))
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), corsMiddleware(cfg.HTTP.AllowedOrigins))

	userCache := memory.NewUserCacheRepository()
	seedRelationshipRepository := memory.NewRelationshipRepository()
	githubClient := githubapi.NewClient(cfg.GitHub)
	relationshipRepository := memory.NewAggregatedRelationshipRepository(seedRelationshipRepository, githubClient)
	bootstrapService := service.NewBootstrapService(cfg.GitHub)
	userService := service.NewUserService(userCache, githubClient)
	pathService := service.NewPathService(relationshipRepository, userService, cfg.Graph.MaxDepth)

	healthHandler := handlerhttp.NewHealthHandler()
	bootstrapHandler := handlerhttp.NewBootstrapHandler(bootstrapService)
	userHandler := handlerhttp.NewUserHandler(userService)
	graphHandler := handlerhttp.NewGraphHandler(pathService)

	router.Register(engine, healthHandler, bootstrapHandler, userHandler, graphHandler)

	return &Application{config: cfg, engine: engine}, nil
}

func (a *Application) Run() error {
	server := &http.Server{
		Addr:         ":" + a.config.App.Port,
		Handler:      a.engine,
		ReadTimeout:  time.Duration(a.config.App.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(a.config.App.WriteTimeoutSeconds) * time.Second,
	}

	return server.ListenAndServe()
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowedSet := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			allowedSet[trimmed] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && originAllowed(origin, allowedSet) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func originAllowed(origin string, allowedSet map[string]struct{}) bool {
	if _, ok := allowedSet["*"]; ok {
		return true
	}

	_, ok := allowedSet[origin]
	return ok
}

func resolveGinMode(environment string) string {
	switch strings.ToLower(strings.TrimSpace(environment)) {
	case "production":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
