package api_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *HttpServer) generateBenchRouter(router *gin.Engine) {
	benchRouter := router.Group("/bench")

	benchRouter.GET("/vllm/start", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	benchRouter.GET("/vllm/stop", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	benchRouter.GET("/vllm/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	benchRouter.GET("/vllm/results", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
