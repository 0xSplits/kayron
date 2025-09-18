package container

// Active defines this worker handler to always be executed.
func (c *Container) Active() bool {
	return true
}
