package reference

type Struct struct {
	// Desired is the release tag or commit sha of this artifact, depending on the
	// specified deployment strategy.
	Desired string
}
