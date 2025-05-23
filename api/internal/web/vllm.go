package api_http

import (
	"context"
	"io"
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
		follow := c.Query("follow")
		if follow == "true" {
			ch := s.vllm.GetLogCh()
			c.Stream(func(w io.Writer) bool {
				line, ok := <-ch
				if !ok {
					return false
				}
				_, err := w.Write([]byte(line))
				return err == nil
			})
		} else {
			logs := s.vllm.GetLogsArchive()
			c.String(http.StatusOK, strings.Join(logs, "\n"))
		}
	})
}
