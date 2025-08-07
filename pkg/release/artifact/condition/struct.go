package condition

type Struct struct {
	// Success indicates whether the desired state requirements for any given
	// scheduler are fulfilled. E.g. the container image tag of any given service
	// release must be pushed to the configured container registry before its
	// desired service release can be rolled out.
	Success bool
}
