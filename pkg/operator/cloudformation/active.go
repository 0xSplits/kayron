package cloudformation

// Active determines whether this worker handler ought to be executed based on
// the underlying policy implementation.
func (c *CloudFormation) Active() bool {
	return c.pol.Update()
}
