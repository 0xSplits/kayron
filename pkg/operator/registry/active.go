package registry

// Active defines this worker handler to always be executed.
func (r *Registry) Active() bool {
	return true
}
