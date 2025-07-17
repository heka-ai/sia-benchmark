package scaleway

import (
	"errors"
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
	for _, server := range servers {
		logger.Debug().Msgf("Server ID: %s, Name: %s, Tags: %s", server.ID, server.Name, server.Tags)

		err = c.deleteServer(*server)
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
	for _, volume := range volumes {
		logger.Debug().Msgf("Volume ID: %s, Name: %s, Tags: %s", volume.ID, volume.Name, volume.Tags)

		err = c.deleteVolume(*volume)
		if err != nil {
			return err
		}
	}

	return nil
}
