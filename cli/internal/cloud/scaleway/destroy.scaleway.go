package scaleway

import (
	"errors"
)

func (c *ScalewayClient) Destroy() error {
	logger.Info().Msg("Mocking Scaleway Destroy")

	err := c.DeleteInstance()
	if err != nil {
		logger.Error().Err(err).Msg("Error while deleting the instance")
		return err
	}

	return nil
}

func (c *ScalewayClient) DeleteInstance() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	servers, err := c.GetAllBenchmarkServers()

	if err != nil {
		return err
	}

	if len(servers) == 0 {
		return errors.New("no server found for this benchmark")
	}

	logger.Debug().Msgf("Found %d servers for this benchmark", len(servers))
	for _, server := range servers {
		logger.Debug().Msgf("Server ID: %s, Name: %s, Tags: %s", server.ID, server.Name, server.Tags)
	}

	for _, server := range servers {
		err = c.deleteServer(server.ID)

		if err != nil {
			return err
		}

		logger.Info().Msgf("Server %s deleted", server.ID)

		for _, volume := range server.Volumes {

			err = c.deleteVolume(volume.ID)

			if err != nil {
				return err
			}

			logger.Info().Msgf("Volume %s deleted", volume.ID)
		}

	}

	return nil
}
