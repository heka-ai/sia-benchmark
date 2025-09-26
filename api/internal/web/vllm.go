package api_http

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *HttpServer) generateVLLMRouter(router *gin.Engine) {
	vllmRouter := router.Group("/vllm")

	vllmRouter.GET("/start", func(c *gin.Context) {
		// Use port from config, with optional query parameter override
		port := s.config.GetConfig().BenchmarkConfig.EnginePort
		if p := c.Query("port"); p != "" {
			if v, err := strconv.Atoi(p); err == nil && v > 0 && v < 65536 {
				port = v
			}
		}

		err := s.vllm.StartWithPort(context.Background(), port)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to start VLLM")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	vllmRouter.GET("/stop", func(c *gin.Context) {
		err := s.vllm.Stop(context.Background())
		if err != nil {
			logger.Error().Err(err).Msg("Failed to stop VLLM")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	vllmRouter.GET("/logs", func(c *gin.Context) {
		logs := s.vllm.GetLogsArchive()
		c.String(http.StatusOK, strings.Join(logs, "\n"))
	})
}
