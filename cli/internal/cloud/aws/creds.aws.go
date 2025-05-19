package aws

func (c *AWSClient) ValidateCredentials() error {
	return c.validateCredentials()
}
