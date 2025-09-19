package policy

// Active defines this worker handler to always be executed.
func (p *Policy) Active() bool {
	return true
}
