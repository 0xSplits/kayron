package infrastructure

// Active determines whether this worker handler ought to be executed based on
// the underlying policy implementation.
func (i *Infrastructure) Active() bool {
	return i.pol.Update()
}
