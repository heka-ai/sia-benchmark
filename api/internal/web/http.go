package api_http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	apiConfig "github.com/heka-ai/benchmark-api/internal/config"
	"github.com/heka-ai/benchmark-api/internal/log"
	"github.com/heka-ai/benchmark-api/pkg/benchmark"
	"github.com/heka-ai/benchmark-api/pkg/vllm"
	"go.uber.org/fx"
)

var logger = log.GetLogger("http")

type HttpServer struct {
	router *gin.Engine

	vllm      *vllm.VLLM
	benchmark *benchmark.Benchmark
	config    *apiConfig.APIConfig
}

var HttpModule = fx.Module("http",
	fx.Provide(NewHttpServer),
)

func NewHttpServer(lc fx.Lifecycle, vllm *vllm.VLLM, benchmark *benchmark.Benchmark, config *apiConfig.APIConfig) *HttpServer {
	server := &HttpServer{
		vllm:      vllm,
		benchmark: benchmark,
		config:    config,
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
		cfg := s.config.GetConfig()
		c.JSON(http.StatusOK, gin.H{
			"status":           "ok",
			"provider":         cfg.GeneralConfig.Provider,
			"inference_engine": cfg.GeneralConfig.InferenceEngine,
			"bench_id":         cfg.GeneralConfig.BenchmarkID,
			"model":            cfg.VLLMConfig.Model,
		})
	})

	s.generateSetupRouter(router)
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
