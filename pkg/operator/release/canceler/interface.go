package canceler

// Interface defines a network abstraction that decides whether any given
// reconciliation loop should be cancelled early on.
type Interface interface {
	// Cancel returns true if the current reconciliation loop should be cancelled.
	Cancel() (bool, error)
}
