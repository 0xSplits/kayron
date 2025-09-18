package infrastructure

// Active defines this worker handler to always be executed.
func (i *Infrastructure) Active() bool {
	return true
}
