package api_http

import (
	"bufio"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/heka-ai/benchmark-api/scripts"
)

var setupScriptFiles = map[string]string{
	"llm":       "gpu_install.sh",
	"benchmark": "cpu_install.sh",
}

func (s *HttpServer) generateSetupRouter(router *gin.Engine) {
	router.POST("/setup", func(c *gin.Context) {
		setupType := c.Query("type")
		scriptName, ok := setupScriptFiles[setupType]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type: must be 'llm' or 'benchmark'"})
			return
		}

		logger.Info().Str("type", setupType).Msg("Running setup")

		scriptData, err := scripts.FS.ReadFile(scriptName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read embedded script: " + err.Error()})
			return
		}

		tmpScript, err := os.CreateTemp("", "bench-setup-*.sh")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer os.Remove(tmpScript.Name())

		if _, err := tmpScript.Write(scriptData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tmpScript.Close()

		if setupType == "benchmark" {
			reqData, _ := scripts.FS.ReadFile("requirements.txt")
			if len(reqData) > 0 {
				os.WriteFile("/tmp/bench-requirements.txt", reqData, 0644)
			}
		}

		cmd := exec.Command("bash", tmpScript.Name())

		stderr, err := cmd.StderrPipe()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := cmd.Start(); err != nil {
			logger.Error().Err(err).Msg("Failed to start setup script")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				logger.Info().Str("setup", setupType).Msg(scanner.Text())
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				logger.Warn().Str("setup", setupType).Msg(scanner.Text())
			}
		}()

		if err := cmd.Wait(); err != nil {
			logger.Error().Err(err).Msg("Setup script failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "setup failed: " + err.Error()})
			return
		}

		logger.Info().Str("type", setupType).Msg("Setup completed successfully")
		c.JSON(http.StatusOK, gin.H{"status": "ok", "type": setupType})
	})
}
