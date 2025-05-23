package api_http

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *HttpServer) generateVLLMRouter(router *gin.Engine) {
	vllmRouter := router.Group("/vllm")

	vllmRouter.GET("/start", func(c *gin.Context) {
		err := s.vllm.Start(context.Background())
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
