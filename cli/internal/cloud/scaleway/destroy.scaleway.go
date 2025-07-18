package scaleway

import (
	"errors"
	"sync"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func (c *ScalewayClient) Destroy() error {
	logger.Info().Msg("Scaleway Destroy")

	err := c.DeleteServers()
	if err != nil {
		logger.Error().Err(err).Msg("Error while deleting the servers")
		return err
	}

	err = c.DeleteVolumes()
	if err != nil {
		logger.Error().Err(err).Msg("Error while deleting the volumes")
		return err
	}

	return nil
}

func (c *ScalewayClient) DeleteServers() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	servers, err := c.GetAllBenchmarkServers()
	if err != nil {
		return err
	}

	logger.Debug().Msgf("Found %d servers for this benchmark", len(servers))
	var wg sync.WaitGroup
	errChan := make(chan error, len(servers))

	for _, server := range servers {
		wg.Add(1)
		go func(srv instance.Server) {
			defer wg.Done()
			logger.Debug().Msgf("Server ID: %s, Name: %s, Tags: %s", srv.ID, srv.Name, srv.Tags)

			if err := c.deleteServer(srv); err != nil {
				errChan <- err
			}
		}(*server)
	}

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

func (c *ScalewayClient) DeleteVolumes() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	volumes, err := c.GetAllBenchmarkVolumes()
	if err != nil {
		return err
	}

	logger.Debug().Msgf("Found %d volumes for this benchmark", len(volumes))
	var wg sync.WaitGroup
	errChan := make(chan error, len(volumes))

	for _, volume := range volumes {
		wg.Add(1)
		go func(vol instance.Volume) {
			defer wg.Done()
			logger.Debug().Msgf("Volume ID: %s, Name: %s, Tags: %s", vol.ID, vol.Name, vol.Tags)

			if err := c.deleteVolume(vol); err != nil {
				errChan <- err
			}
		}(*volume)
	}

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
