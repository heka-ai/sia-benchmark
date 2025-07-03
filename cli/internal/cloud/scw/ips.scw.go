package scw

func (c *SCWClient) GetLLMInstanceIP() (string, error) {
	logger.Info().Msg("Mocking Scaleway LLM IP")

	return "127.0.0.1", nil
}

func (c *SCWClient) GetBenchInstanceIP() (string, error) {
	logger.Info().Msg("Mocking Scaleway Bench IP")

	return "127.0.0.2", nil
}
