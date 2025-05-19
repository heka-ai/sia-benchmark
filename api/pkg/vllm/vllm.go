package vllm

import (
	"bufio"
	"context"
	"os"
	"os/exec"

	apiConfig "github.com/heka-ai/benchmark-api/internal/config"
	"github.com/heka-ai/benchmark-api/internal/log"
	"go.uber.org/fx"
)

var logger = log.GetLogger("vllm")

var PATH_TO_VLLM = "/usr/local/bin/vllm"

var VLLMModule = fx.Module("vllm",
	fx.Provide(NewVLLM),
)

type VLLM struct {
	args    []string
	cmd     *exec.Cmd
	doneCh  chan struct{}
	waitCh  chan struct{}
	running int64
	config  *apiConfig.APIConfig
}

func NewVLLM(lc fx.Lifecycle, config *apiConfig.APIConfig) *VLLM {
	vllm := &VLLM{
		args:    []string{},
		doneCh:  make(chan struct{}),
		waitCh:  make(chan struct{}),
		running: 0,
		config:  config,
	}

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return vllm.Stop(ctx)
	}))

	return vllm
}

func (v *VLLM) Start(ctx context.Context) error {
	logger.Info().Str("model", v.config.GetConfig().VLLMConfig.Model).Str("token", v.config.GetConfig().BenchmarkConfig.Token).Msg("Starting the VLLM service")

	localArgs := []string{"serve", v.config.GetConfig().VLLMConfig.Model}

	v.cmd = exec.CommandContext(ctx, PATH_TO_VLLM, localArgs...)

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
			logger.Info().Msg(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			logger.Warn().Msg(scanner.Text())
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

	return v.cmd.Process.Kill()
}
