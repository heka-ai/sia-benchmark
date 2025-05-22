package bench

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/heka-ai/benchmark-cli/pkg/results"
)

// This client follow the Sia Benchmark API
// it is used to interact with the differents instances
type Client struct {
	APIKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	httpClient := &http.Client{}
	return &Client{
		APIKey:     apiKey,
		httpClient: httpClient,
	}
}

const (
	waitInterval  = 1 * time.Second
	maxIterations = 100
)

func (c *Client) WaitForInstances(benchIP, llmIP string) error {
	cpuDone := false
	llmDone := false

	for i := 0; i < maxIterations && (!cpuDone || !llmDone); i++ {
		err := c.HealthCheck(benchIP)
		if err == nil {
			cpuDone = true
		}

		err = c.HealthCheck(llmIP)
		if err == nil {
			llmDone = true
		}

		time.Sleep(waitInterval)
	}

	return fmt.Errorf("instances are not running after %d iterations", maxIterations)
}

// deploy the model on the instance
// will also upload the config to the instance
func (c *Client) Deploy(ip string, engine string) error {
	// config to string
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:8001/%s/start", ip, engine), nil)
	if err != nil {
		return err
	}

	request.Header.Add("X-API-Key", c.APIKey)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deploy: %s", resp.Status)
	}

	return nil
}

func (c *Client) HealthCheck(ip string) error {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:8001/health", ip), nil)
	if err != nil {
		return err
	}

	request.Header.Add("X-API-Key", c.APIKey)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to health check: %s", resp.Status)
	}

	return nil
}

func (c *Client) WaitForLLM(ip string) error {
	done := false

	for i := 0; i < maxIterations && !done; i++ {
		done, _ := c.ModelStatus(ip)

		if done {
			return nil
		}

		time.Sleep(waitInterval)
	}

	return fmt.Errorf("LLM is not ready after %d iterations", maxIterations)
}

func (c *Client) ModelStatus(ip string) (bool, error) {
	res, err := http.Get(fmt.Sprintf("http://%s:8000/v1/models", ip))
	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to get model status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// check if the body is a json and the vllm model is ready
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("failed to parse model status response: %v", err)
	}

	if status, ok := result["status"].(string); ok && status == "ready" {
		return true, nil
	}

	return false, nil
}

func (c *Client) RunBenchmark(ip string, engineType string) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s:8001/%s/start", ip, engineType), nil)
	if err != nil {
		return err
	}

	request.Header.Add("X-API-Key", c.APIKey)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to run benchmark: %s", resp.Status)
	}

	return nil
}

func (c *Client) GetResults(ip string, engineType string) (*results.Results, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:8001/%s/results", ip, engineType), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-API-Key", c.APIKey)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get results: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results results.Results
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse results: %v", err)
	}

	// validate
	validator := validator.New()
	if err := validator.Struct(results); err != nil {
		return nil, fmt.Errorf("failed to validate results: %v", err)
	}

	return &results, nil
}
