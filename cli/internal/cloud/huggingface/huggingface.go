package huggingface

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/heka-ai/benchmark-cli/internal/bench"
	"github.com/heka-ai/benchmark-cli/internal/cloud"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

var logger = log.GetLogger("huggingface")

type HFClient struct {
	cloud.Cloud

	cli     *bench.Client
	config  *config.Config
	wasInit bool
}

func NewClient(cfg *config.Config) *HFClient {
	return &HFClient{
		cli:     bench.NewClient(cfg.APIKey),
		config:  cfg,
		wasInit: false,
	}
}

func (c *HFClient) Init() cloud.Cloud {
	c.wasInit = true
	return c
}

func (c *HFClient) ValidateCredentials() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	// Skip validation when using random dataset (no HF access required)
	if c.config != nil && c.config.BenchmarkConfig != nil && c.config.BenchmarkConfig.DatasetName == "random" {
		logger.Info().Msg("Skipping HuggingFace token validation (dataset_name is 'random')")
		return nil
	}

	token := ""
	if c.config != nil && c.config.BenchmarkConfig != nil {
		token = c.config.BenchmarkConfig.Token
	}
	if token == "" {
		return errors.New("missing HuggingFace token (benchmark.token or HF_TOKEN or HUGGINGFACEHUB_API_TOKEN or HF_API_TOKEN)")
	}

	// 1) Validate token via whoami
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://huggingface.co/api/whoami-v2", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sia-benchmark/cli")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info().Msg("OK - HuggingFace token is valid")
		// continue to model access check
	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.New("invalid HuggingFace token")
	default:
		return fmt.Errorf("unexpected response validating token: %s", resp.Status)
	}

	// 2) Ensure access to configured model
	if c.config == nil || c.config.VLLMConfig == nil || c.config.VLLMConfig.Model == "" {
		return errors.New("missing vllm.model in config")
	}
	model := c.config.VLLMConfig.Model

	modelReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://huggingface.co/api/models/"+model, nil)
	if err != nil {
		return err
	}
	modelReq.Header.Set("Authorization", "Bearer "+token)
	modelReq.Header.Set("Accept", "application/json")
	modelReq.Header.Set("User-Agent", "sia-benchmark/cli")

	modelResp, err := client.Do(modelReq)
	if err != nil {
		return err
	}
	defer modelResp.Body.Close()

	switch modelResp.StatusCode {
	case http.StatusOK:
		logger.Info().Str("model", model).Msg("OK - Access to model confirmed")
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return fmt.Errorf("no access to model '%s' with provided token", model)
	case http.StatusNotFound:
		return fmt.Errorf("model '%s' not found or not accessible", model)
	default:
		return fmt.Errorf("unexpected response checking model '%s': %s", model, modelResp.Status)
	}
}

func (c *HFClient) Create() error { return errors.New("not supported for huggingface") }
func (c *HFClient) CreateLLMInstance() error { return errors.New("not supported for huggingface") }
func (c *HFClient) CreateBenchInstance() error { return errors.New("not supported for huggingface") }
func (c *HFClient) Destroy() error { return errors.New("not supported for huggingface") }
func (c *HFClient) GetLLMInstanceIP() (string, error) { return "", errors.New("not supported for huggingface") }
func (c *HFClient) GetBenchInstanceIP() (string, error) { return "", errors.New("not supported for huggingface") }
