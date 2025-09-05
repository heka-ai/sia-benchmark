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
	return v.StartWithPort(ctx, 8000)
}

// StartWithPort launches vLLM specifying the serving port. It appends
// `--port <port>` to the generated vllm arguments.
func (v *VLLM) StartWithPort(ctx context.Context, port int) error {
	logger.Info().Str("model", v.config.GetConfig().VLLMConfig.Model).Str("token", v.config.GetConfig().BenchmarkConfig.Token).Msg("Starting the VLLM service")

	localArgs, err := cliConfig.GenerateVLLMCommand(v.config.GetConfig().VLLMConfig)
	if err != nil {
		return err
	}

	// Inject the port flag for vllm serve.
	localArgs = append(localArgs, "--port", fmt.Sprintf("%d", port))

	logger.Info().Str("command", "vllm "+strings.Join(localArgs, " ")).Msg("Launching VLLM with the following command")

	// Resolve the path to the vllm binary to ensure vllm is installed on host machine.
	logger.Info().Str("PATH", os.Getenv("PATH")).Msg("Environment PATH before resolving vllm")
	bin, err := resolveVLLMBinary()
	if err != nil {
		return err
	}

	v.cmd = exec.CommandContext(ctx, bin, localArgs...)
	v.cmd.Env = append(os.Environ(), "HF_TOKEN="+v.config.GetConfig().BenchmarkConfig.Token)

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
		logger.Info().Msg("VLLM process not started; nothing to stop")
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

// resolveVLLMBinary determines the path to the vllm executable.
func resolveVLLMBinary() (string, error) {
	if p, err := exec.LookPath("vllm"); err == nil {
		return p, nil
	}
	// Fallback to common installation locations
	candidates := []string{
		"/opt/pytorch/bin/vllm",
		"/usr/local/bin/vllm",
		"/usr/bin/vllm",
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			if st.Mode()&0111 != 0 {
				return c, nil
			}
		}
	}
	return "", fmt.Errorf("vllm not installed or not in PATH (PATH=%s)", os.Getenv("PATH"))
}
