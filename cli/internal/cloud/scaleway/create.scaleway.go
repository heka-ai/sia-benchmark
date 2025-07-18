package scaleway

import (
	"errors"
	"sync"
)

func (c *ScalewayClient) Create() error {
	logger.Info().Msg("Scaleway Create")

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.CreateBenchmarkClient(); err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.CreateBenchmarkServer(); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ScalewayClient) CreateBenchmarkClient() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	err := c.createServer()
	if err != nil {
		return err
	}

	return nil
}

func (c *ScalewayClient) CreateBenchmarkServer() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	err := c.createServer()
	if err != nil {
		return err
	}

	return nil
}
