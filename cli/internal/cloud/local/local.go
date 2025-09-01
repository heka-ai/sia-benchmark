package local

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/heka-ai/benchmark-cli/internal/bench"
	"github.com/heka-ai/benchmark-cli/internal/cloud"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

var logger = log.GetLogger("local")

type LocalClient struct {
	cloud.Cloud

	cli     *bench.Client
	config  *config.Config
	wasInit bool
}

func NewClient(cfg *config.Config) *LocalClient {
	return &LocalClient{
		cli:     bench.NewClient(cfg.APIKey),
		config:  cfg,
		wasInit: false,
	}
}

func (c *LocalClient) Init() cloud.Cloud {
	c.wasInit = true
	return c
}

func (c *LocalClient) ValidateCredentials() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	// Try nvidia-smi first
	nameMem, nvidiaErr := queryNvidiaSMI()
	if nvidiaErr == nil && len(nameMem) > 0 {
		for i, nm := range nameMem {
			logger.Info().Int("index", i).Str("name", nm.name).Str("memory", nm.memory).Msg("Detected NVIDIA GPU via nvidia-smi")
		}
		return nil
	}

	// Fallback to tegrastats (Jetson)
	nameMem, tegraErr := queryTegrastats()
	if tegraErr == nil && len(nameMem) > 0 {
		for i, nm := range nameMem {
			logger.Info().Int("index", i).Str("name", nm.name).Str("memory", nm.memory).Msg("Detected NVIDIA GPU via tegrastats")
		}
		return nil
	}

	// Log detailed errors
	if nvidiaErr != nil {
		logger.Error().Err(nvidiaErr).Msg("nvidia-smi failed or is unavailable")
	}
	if tegraErr != nil {
		logger.Error().Err(tegraErr).Msg("tegrastats failed or is unavailable")
	}

	logger.Error().Msg("Neither nvidia-smi nor tegrastats worked. GPU detection failed.")
	return errors.New("gpu detection failed: neither nvidia-smi nor tegrastats worked")
}

type gpuInfo struct {
	name   string
	memory string
}

func queryNvidiaSMI() ([]gpuInfo, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	var infos []gpuInfo
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		// Expect: "NVIDIA A100-SXM4-40GB, 40536 MiB"
		parts := strings.Split(l, ",")
		if len(parts) < 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		mem := strings.TrimSpace(parts[1])
		infos = append(infos, gpuInfo{name: name, memory: mem})
	}
	if len(infos) == 0 {
		return nil, errors.New("no gpu returned by nvidia-smi")
	}
	return infos, nil
}

func queryTegrastats() ([]gpuInfo, error) {
	if _, err := exec.LookPath("tegrastats"); err != nil {
		return nil, errors.New("tegrastats command not found")
	}

	// Capture a single sample line
	cmd := exec.Command("sh", "-c", "tegrastats --interval 200 | head -n 1")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	line := strings.TrimSpace(out.String())
	if line == "" {
		return nil, errors.New("empty tegrastats output")
	}

	// GPU presence indicator (Jetson) is GR3D_FREQ
	if !strings.Contains(line, "GR3D_FREQ") {
		return nil, errors.New("tegrastats output missing GR3D_FREQ (GPU)")
	}

	// Jetson uses unified/shared memory. Use total system RAM as an approximation.
	memTotal := ""
	ramTotalRe := regexp.MustCompile(`RAM\s+\d+/(\d+)MB`)
	if m := ramTotalRe.FindStringSubmatch(line); len(m) == 2 {
		memTotal = m[1] + " MB (shared)"
	}

	name := "Jetson-GPU"
	if b, err := os.ReadFile("/proc/device-tree/model"); err == nil {
		name = strings.TrimSpace(string(b))
	}

	if memTotal == "" {
		memTotal = "unknown"
	}

	return []gpuInfo{{name: name, memory: memTotal}}, nil
}

func (c *LocalClient) Create() error {
	logger.Warn().Msg("Local provider: nothing to create (no-op)")
	return nil
}
func (c *LocalClient) CreateLLMInstance() error {
	logger.Warn().Msg("Local provider: LLM instance creation is a no-op")
	return nil
}
func (c *LocalClient) CreateBenchInstance() error {
	logger.Warn().Msg("Local provider: Bench instance creation is a no-op")
	return nil
}
func (c *LocalClient) Destroy() error {
	logger.Warn().Msg("Local provider: nothing to destroy (no-op)")
	return nil
}

// For local, both endpoints are localhost
func (c *LocalClient) GetLLMInstanceIP() (string, error) { return "127.0.0.1", nil }
func (c *LocalClient) GetBenchInstanceIP() (string, error) { return "127.0.0.1", nil }
