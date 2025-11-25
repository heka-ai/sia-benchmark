package vllm

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"go.uber.org/fx"

	cliConfig "github.com/heka-ai/benchmark-cli/pkg/config"

	apiConfig "github.com/heka-ai/benchmark-api/internal/config"
	"github.com/heka-ai/benchmark-api/internal/log"
)

var logger = log.GetLogger("vllm")

var VLLMModule = fx.Module("vllm",
	fx.Provide(NewVLLM),
)

type VLLM struct {
	args        []string
	cmd         *exec.Cmd
	doneCh      chan struct{}
	waitCh      chan struct{}
	running     int64
	logsArchive []string
	config      *apiConfig.APIConfig
}

func NewVLLM(lc fx.Lifecycle, config *apiConfig.APIConfig) *VLLM {
	vllm := &VLLM{
		args:        []string{},
		doneCh:      make(chan struct{}),
		waitCh:      make(chan struct{}),
		logsArchive: []string{},
		running:     0,
		config:      config,
	}

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return vllm.Stop(ctx)
	}))

	return vllm
}

func (v *VLLM) GetLogsArchive() []string {
	return v.logsArchive
}

func (v *VLLM) Start(ctx context.Context) error {
	cfg := v.config.GetConfig()
	hfToken := ""
	if cfg.BenchmarkConfig != nil && strings.TrimSpace(cfg.BenchmarkConfig.Token) != "" {
		hfToken = strings.TrimSpace(cfg.BenchmarkConfig.Token)
	} else {
		hfToken = strings.TrimSpace(os.Getenv("HF_TOKEN"))
	}
	if hfToken == "" {
		logger.Warn().Msg("Hugging Face token is empty; set benchmark.token or HF_TOKEN before deploy for gated/private models")
	}

	logger.Info().Str("model", cfg.VLLMConfig.Model).Msg("Starting the VLLM service")

	localArgs, err := cliConfig.GenerateVLLMCommand(cfg.VLLMConfig)
	if err != nil {
		return err
	}

	logger.Info().Str("command", "uv run vllm "+strings.Join(localArgs, " ")).Msg("Launching VLLM with the following command")

	logger.Info().Str("command", "uv run vllm "+strings.Join(localArgs, " ")).Msg("Launching VLLM with the following command")

	// Resolve the path to uv which orchestrates the Python environment for vllm.
	logger.Info().Str("PATH", os.Getenv("PATH")).Msg("Environment PATH before resolving uv")
	uvBin, err := exec.LookPath("uv")
	if err != nil {
		return fmt.Errorf("uv not found in PATH: please install uv or adjust PATH")
	}

	uvArgs := append([]string{"run", "vllm"}, localArgs...)

	v.cmd = exec.CommandContext(ctx, uvBin, uvArgs...)
	v.cmd.Env = append(os.Environ(),
		"HF_TOKEN="+hfToken,
		"HF_HUB_ENABLE_HF_TRANSFER=1",
	)

	stdout, err := v.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := v.cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			v.logsArchive = append(v.logsArchive, line)
			logger.Info().Msg(line)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			v.logsArchive = append(v.logsArchive, line)
			logger.Warn().Msg(line)
		}
	}()

	if err := v.cmd.Start(); err != nil {
		return err
	}

	go func() {
		v.cmd.Wait()
		close(v.doneCh)
	}()

	return nil
}

func (v *VLLM) Stop(ctx context.Context) error {
	logger.Info().Msg("Stopping VLLM")

	// If Start was never called or the process was not created, there is nothing to stop.
	if v == nil || v.cmd == nil || v.cmd.Process == nil {
		logger.Error().Msg("VLLM process not started; nothing to stop")
		return nil
	}

	// Attempt a graceful shutdown first: send SIGTERM and wait for the process to exit.
	_ = v.cmd.Process.Signal(syscall.SIGTERM)

	select {
	case <-v.doneCh:
		return nil
	case <-time.After(5 * time.Second):
		// Force terminate if it did not exit in time.
		return v.cmd.Process.Kill()
	}
}
