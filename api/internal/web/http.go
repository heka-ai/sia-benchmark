package api_http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	apiConfig "github.com/heka-ai/benchmark-api/internal/config"
	"github.com/heka-ai/benchmark-api/internal/log"
	"github.com/heka-ai/benchmark-api/pkg/vllm"
	"go.uber.org/fx"
)

var logger = log.GetLogger("http")

type HttpServer struct {
	router *gin.Engine

	vllm   *vllm.VLLM
	config *apiConfig.APIConfig
}

var HttpModule = fx.Module("http",
	fx.Provide(NewHttpServer),
)

func NewHttpServer(lc fx.Lifecycle, vllm *vllm.VLLM, config *apiConfig.APIConfig) *HttpServer {
	server := &HttpServer{
		vllm:   vllm,
		config: config,
	}
	server.router = server.createRouter()

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		return server.Start(ctx)
	}))

	return server
}

func (s *HttpServer) createRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	gin.DefaultWriter = log.GetMainLogger().With().Str("level", "info").Str("module", "http").Logger()
	gin.DefaultErrorWriter = log.GetMainLogger().With().Str("level", "error").Str("module", "http").Logger()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "provider": s.config.GetConfig().Provider, "inference_engine": s.config.GetConfig().InferenceEngine, "bench_id": s.config.GetConfig().BenchID, "model": s.config.GetConfig().VLLMConfig.Model})
	})

	// generate the vllm routes
	s.generateVLLMRouter(router)
	s.generateBenchRouter(router)

	return router
}

func (s *HttpServer) Start(ctx context.Context) error {
	logger.Info().Str("address", ":8001").Msg("Starting the HTTP server")

	server := &http.Server{
		Addr:    ":8001",
		Handler: s.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	return nil
}
