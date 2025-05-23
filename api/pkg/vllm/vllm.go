package vllm

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"strings"

	apiConfig "github.com/heka-ai/benchmark-api/internal/config"
	"github.com/heka-ai/benchmark-api/internal/log"
	cliConfig "github.com/heka-ai/benchmark-cli/pkg/config"
	"go.uber.org/fx"
)

var logger = log.GetLogger("vllm")

var PATH_TO_VLLM = "/usr/local/bin/vllm"

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
	logger.Info().Str("model", v.config.GetConfig().VLLMConfig.Model).Str("token", v.config.GetConfig().BenchmarkConfig.Token).Msg("Starting the VLLM service")

	localArgs, err := cliConfig.GenerateVLLMCommand(v.config.GetConfig().VLLMConfig)
	if err != nil {
		return err
	}

	logger.Info().Str("command", "vllm "+strings.Join(localArgs, " ")).Msg("Launching VLLM with the following command")

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

	return v.cmd.Process.Kill()
}
