package suspend

// Suspend disables any further reconciliation of this service indefinitely.
type Suspend bool

func (s Suspend) Empty() bool {
	return !bool(s)
}

func (s Suspend) Verify() error {
	// TODO
	return nil
}
