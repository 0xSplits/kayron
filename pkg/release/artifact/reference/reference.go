package reference

type Struct struct {
	// Current is the image tag currently deployed within the underlying
	// infrastructure provider, e.g. CloudFormation.
	Current string

	// Desired is either the release tag or commit sha for this artifact,
	// depending on the specified deployment startegy.
	Desired string
}
