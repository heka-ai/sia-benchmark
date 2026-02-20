package api_http

import (
	"bufio"
	"net/http"
	"os"
	"os/exec"
	"sync"

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
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lines := make(chan string)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}()
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				lines <- "stderr: " + scanner.Text()
			}
		}()
		go func() { wg.Wait(); close(lines) }()

		c.Header("Content-Type", "text/plain; charset=utf-8")
		flusher, _ := c.Writer.(http.Flusher)
		for line := range lines {
			c.Writer.Write([]byte(line + "\n"))
			if flusher != nil {
				flusher.Flush()
			}
		}

		if err := cmd.Wait(); err != nil {
			c.Writer.Write([]byte("error: " + err.Error() + "\n"))
			if flusher != nil {
				flusher.Flush()
			}
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
}
