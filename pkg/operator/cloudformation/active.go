package cloudformation

// Active defines this worker handler to always be executed.
func (c *CloudFormation) Active() bool {
	return true
}
