package scheduler

type Struct struct {
	// Current is either the container image tag of any given service or the
	// template version parameter for any given infrastructure currently deployed
	// within the configured runtime.
	Current string
}
