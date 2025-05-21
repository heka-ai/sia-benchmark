package benchmark

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

var PATH_TO_PYTHON = "/opt/pytorch/bin/python3"

var logger = log.GetLogger("benchmark")

type Benchmark struct {
	args    []string
	cmd     *exec.Cmd
	doneCh  chan struct{}
	waitCh  chan struct{}
	running int64

	config *apiConfig.APIConfig
}

func NewBenchmark(lc fx.Lifecycle, config *apiConfig.APIConfig) *Benchmark {
	benchmark := &Benchmark{
		args:    []string{},
		doneCh:  make(chan struct{}),
		waitCh:  make(chan struct{}),
		running: 0,
		config:  config,
	}

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return benchmark.Stop()
	}))

	return benchmark
}

func (b *Benchmark) Start(ip string) error {
	logger.Info().Str("command", "benchmark "+strings.Join(b.args, " ")).Msg("Starting benchmark")

	localArgs, err := cliConfig.GenerateBenchmarkCommand(b.config.GetConfig(), ip)
	if err != nil {
		return err
	}

	b.cmd = exec.CommandContext(context.Background(), PATH_TO_PYTHON, localArgs...)
	b.cmd.Env = append(os.Environ(), "HF_TOKEN="+b.config.GetConfig().BenchmarkConfig.Token)

	stdout, err := b.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := b.cmd.StderrPipe()
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

	if err := b.cmd.Start(); err != nil {
		return err
	}

	go func() {
		b.cmd.Wait()
		close(b.doneCh)
	}()

	return nil
}

func (b *Benchmark) Stop() error {
	return nil
}
