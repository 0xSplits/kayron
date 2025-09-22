package status

// Active defines this worker handler to always be executed.
func (s *Status) Active() bool {
	return true
}
