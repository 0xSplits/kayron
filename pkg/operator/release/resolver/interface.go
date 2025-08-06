package resolver

// Interface defines a network abstraction that we want to control in order to
// verify the business logic of resolver.Search during unit tests.
type Interface interface {
	// Commit returns the most recent Git commit at the given ref for the
	// underlying Github repository.
	Commit(string) (string, error)

	// Latest returns the most recent Git tag for the latest release of the
	// underlying Github repository.
	Latest() (string, error)
}
