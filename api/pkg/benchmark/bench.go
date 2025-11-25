package benchmark

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"

	cliConfig "github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/heka-ai/benchmark-cli/pkg/results"

	apiConfig "github.com/heka-ai/benchmark-api/internal/config"
	"github.com/heka-ai/benchmark-api/internal/log"
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

	config     *apiConfig.APIConfig
	resultPath string
}

var BenchmarkModule = fx.Module("benchmark",
	fx.Provide(NewBenchmark),
	fx.Invoke(func(b *Benchmark) {}),
)

func (b *Benchmark) GetLogsArchive() []string {
	return b.logsArchive
}

func NewBenchmark(lc fx.Lifecycle, config *apiConfig.APIConfig) *Benchmark {
	benchmark := &Benchmark{
		args:        []string{},
		doneCh:      make(chan struct{}),
		waitCh:      make(chan struct{}),
		logsArchive: []string{},
		running:     0,
		config:      config,
	}

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return benchmark.Stop()
	}))

	return benchmark
}

func (b *Benchmark) Start(ip string) error {
	cfg := b.config.GetConfig()
	port := cfg.BenchmarkConfig.EnginePort
	resultFilename := strings.TrimSpace(cfg.BenchmarkConfig.ResultFilename)
	if resultFilename == "" {
		resultFilename = PATH_TO_RESULTS
	}

	localArgs, err := cliConfig.GenerateBenchmarkCommand(cfg, ip, port)
	if err != nil {
		return err
	}

	b.resultPath = resultFilename
	localArgs = append(localArgs, "--save-result", "--result-filename", resultFilename)

	// Resolve uv binary via PATH (now includes /home/ubuntu/.local/bin)
	p, err := exec.LookPath("uv")
	if err != nil {
		return fmt.Errorf("uv not found in PATH: please install uv or adjust PATH")
	}

	// Prepend "run python3" to the arguments
	uvArgs := append([]string{"run", "python3"}, localArgs...)

	logger.Info().Str("command", p+" "+strings.Join(uvArgs, " ")).Msg("Starting benchmark")

	b.cmd = exec.CommandContext(context.Background(), p, uvArgs...)

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
			b.logsArchive = append(b.logsArchive, line)
			logger.Info().Msg(line)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
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
	path := b.resultPath
	if strings.TrimSpace(path) == "" {
		path = PATH_TO_RESULTS
	}
	file, err := os.Open(path)
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
