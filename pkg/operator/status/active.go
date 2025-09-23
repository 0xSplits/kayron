package status

// Active defines this worker handler to only be executed within the testing
// environment, because we do not allow preview deployments to be injected in
// e.g. staging nor production.
func (s *Status) Active() bool {
	return s.env.Environment == "testing"
}
