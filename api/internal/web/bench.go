package api_http

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BenchStartRequest struct {
	IP string `json:"ip"`
}

func (s *HttpServer) generateBenchRouter(router *gin.Engine) {
	benchRouter := router.Group("/bench")

	benchRouter.POST("/vllm/start", func(c *gin.Context) {
		var req BenchStartRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logger.Info().Str("ip", req.IP).Msg("Starting benchmark")

		err := s.benchmark.Start(req.IP)

		if err != nil {
			logger.Error().Err(err).Msg("Failed to start benchmark")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	benchRouter.GET("/vllm/stop", func(c *gin.Context) {
		err := s.benchmark.Stop()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to stop benchmark")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	benchRouter.GET("/vllm/status", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"status": "not implemented"})
	})

	benchRouter.GET("/vllm/logs", func(c *gin.Context) {
		follow := c.Query("follow")
		if follow == "true" {
			ch := s.benchmark.GetLogCh()
			c.Stream(func(w io.Writer) bool {
				line, ok := <-ch
				if !ok {
					return false
				}
				_, err := w.Write([]byte(line))
				if err != nil {
					logger.Error().Err(err).Msg("Failed to write line")
					return false
				}
				return true
			})
		} else {
			logs := s.benchmark.GetLogsArchive()
			c.String(http.StatusOK, strings.Join(logs, "\n"))
		}
	})

	benchRouter.GET("/vllm/results", func(c *gin.Context) {
		results, err := s.benchmark.GetResult()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get results")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	})
}
