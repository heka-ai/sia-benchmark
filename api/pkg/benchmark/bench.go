package benchmark

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/go-playground/validator/v10"
	apiConfig "github.com/heka-ai/benchmark-api/internal/config"
	"github.com/heka-ai/benchmark-api/internal/log"
	cliConfig "github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/heka-ai/benchmark-cli/pkg/results"
	"go.uber.org/fx"
)

var PATH_TO_PYTHON = "/opt/pytorch/bin/python3"
var PATH_TO_RESULTS = "/home/ubuntu/metrics.json"

var logger = log.GetLogger("benchmark")

type Benchmark struct {
	args    []string
	cmd     *exec.Cmd
	doneCh  chan struct{}
	waitCh  chan struct{}
	running int64

	logsArchive []string
	logCh       chan string

	config *apiConfig.APIConfig
}

var BenchmarkModule = fx.Module("benchmark",
	fx.Provide(NewBenchmark),
	fx.Invoke(func(b *Benchmark) {}),
)

func (b *Benchmark) GetLogsArchive() []string {
	return b.logsArchive
}

func (b *Benchmark) GetLogCh() chan string {
	return b.logCh
}

func NewBenchmark(lc fx.Lifecycle, config *apiConfig.APIConfig) *Benchmark {
	benchmark := &Benchmark{
		args:        []string{},
		doneCh:      make(chan struct{}),
		waitCh:      make(chan struct{}),
		logsArchive: []string{},
		logCh:       make(chan string),
		running:     0,
		config:      config,
	}

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return benchmark.Stop()
	}))

	return benchmark
}

func (b *Benchmark) Start(ip string) error {
	localArgs, err := cliConfig.GenerateBenchmarkCommand(b.config.GetConfig(), ip)
	if err != nil {
		return err
	}

	localArgs = append(localArgs, "--save_result", "true", "--result_filename", PATH_TO_RESULTS)

	logger.Info().Str("command", PATH_TO_PYTHON+" "+strings.Join(localArgs, " ")).Msg("Starting benchmark")

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
			line := scanner.Text()
			b.logCh <- line
			b.logsArchive = append(b.logsArchive, line)
			logger.Info().Msg(line)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			b.logCh <- line
			b.logsArchive = append(b.logsArchive, line)
			logger.Warn().Msg(line)
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

func (b *Benchmark) GetResult() (*results.Results, error) {
	file, err := os.Open(PATH_TO_RESULTS)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var results results.Results
	if err := json.Unmarshal(bytes, &results); err != nil {
		return nil, err
	}

	// validate
	val := validator.New()
	if err := val.Struct(results); err != nil {
		return nil, err
	}

	return &results, nil
}

func (b *Benchmark) Stop() error {
	return nil
}
