package container

type Struct struct {
	// Exists expresses whether the desired artifact reference exists as Docker
	// image inside of the configured Docker registry.
	Exists bool
}
