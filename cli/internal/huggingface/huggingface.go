package huggingface

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

var logger = log.GetLogger("huggingface")

// ValidateHFCredentials validates the HuggingFace token and access to the configured model.
// Lightweight functional API (no cloud abstraction).
func ValidateHFCredentials(cfg *config.Config) error {
	token := ""
	if cfg != nil && cfg.BenchmarkConfig != nil {
		token = cfg.BenchmarkConfig.Token
	}
	// Fallback to environment variables if token is empty
	if token == "" {
		for _, envKey := range []string{"HF_TOKEN", "HUGGINGFACEHUB_API_TOKEN", "HF_API_TOKEN"} {
			if env := os.Getenv(envKey); strings.TrimSpace(env) != "" {
				token = env
				break
			}
		}
	}
	if token == "" {
		return errors.New("missing HuggingFace token (benchmark.token or HF_TOKEN or HUGGINGFACEHUB_API_TOKEN or HF_API_TOKEN)")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// 1) Validate token via whoami
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://huggingface.co/api/whoami-v2", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sia-benchmark/cli")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info().Msg("OK - HuggingFace token is valid")
	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.New("invalid HuggingFace token")
	default:
		return fmt.Errorf("unexpected response validating token: %s", resp.Status)
	}

	// 2) Ensure access to configured model
	if cfg == nil || cfg.VLLMConfig == nil || cfg.VLLMConfig.Model == "" {
		return errors.New("missing vllm.model in config")
	}
	model := cfg.VLLMConfig.Model

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
	case http.StatusUnauthorized, http.StatusForbidden:
		return fmt.Errorf("no access to model '%s' with provided token", model)
	case http.StatusNotFound:
		return fmt.Errorf("model '%s' not found or not accessible", model)
	default:
		return fmt.Errorf("unexpected response checking model '%s': %s", model, modelResp.Status)
	}

	// 3) If dataset_name is 'hf', ensure access to configured dataset
	if cfg.BenchmarkConfig != nil {
		datasetName := cfg.BenchmarkConfig.DatasetName
		if datasetName == "hf" { // skip when 'random', only validate for HF datasets
			datasetId := strings.TrimSpace(cfg.BenchmarkConfig.DatasetPath)
			if datasetId == "" {
				return errors.New("missing benchmark.dataset_path for HuggingFace dataset")
			}

			dsReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://huggingface.co/api/datasets/"+datasetId, nil)
			if err != nil {
				return err
			}
			dsReq.Header.Set("Authorization", "Bearer "+token)
			dsReq.Header.Set("Accept", "application/json")
			dsReq.Header.Set("User-Agent", "sia-benchmark/cli")

			dsResp, err := client.Do(dsReq)
			if err != nil {
				return err
			}
			defer dsResp.Body.Close()

			switch dsResp.StatusCode {
			case http.StatusOK:
				logger.Info().Str("dataset", datasetId).Msg("OK - Access to dataset confirmed")
			case http.StatusUnauthorized, http.StatusForbidden:
				return fmt.Errorf("no access to dataset '%s' with provided token", datasetId)
			case http.StatusNotFound:
				return fmt.Errorf("dataset '%s' not found or not accessible", datasetId)
			default:
				return fmt.Errorf("unexpected response checking dataset '%s': %s", datasetId, dsResp.Status)
			}
		}
	}

	return nil
}
