package resolver

// Interface defines a network abstraction that we want to control in order to
// verify the business logic of resolver.Search during unit tests.
type Interface interface {
	// Exists determines whether the given branch exists for the underlying Github
	// repository.
	Exists(string) (bool, error)

	// Latest returns the most recent Git tag for the latest release of the
	// underlying Github repository.
	Latest() (string, error)
}
