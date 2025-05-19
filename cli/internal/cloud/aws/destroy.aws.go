package aws

func (c *AWSClient) Destroy() error {
	err := c.DeleteInstance()
	if err != nil {
		logger.Error().Err(err).Msg("Error while deleting the instance")
		return err
	}

	return nil
}
