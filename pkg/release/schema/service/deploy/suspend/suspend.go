package suspend

// Bool disables any further reconciliation of this service indefinitely.
type Bool bool

func (s Bool) Empty() bool {
	return !bool(s)
}

func (s Bool) Verify() error {
	// TODO
	return nil
}
