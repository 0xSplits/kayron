package release

// Active defines this worker handler to always be executed.
func (r *Release) Active() bool {
	return true
}
