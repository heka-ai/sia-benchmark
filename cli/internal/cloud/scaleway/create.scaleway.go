package scaleway

import "errors"

func (c *ScalewayClient) Create() error {
	logger.Info().Msg("Scaleway Create")

	err := c.CreateServer()
	if err != nil {
		logger.Error().Err(err).Msg("Error while creating the instance")
		return err
	}

	return nil
}

func (c *ScalewayClient) CreateServer() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	err := c.createServer()
	if err != nil {
		return err
	}

	return nil
}
