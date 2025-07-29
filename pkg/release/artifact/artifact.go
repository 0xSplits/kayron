package artifact

import (
	"github.com/0xSplits/kayron/pkg/release/artifact/container"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
)

// Artifact describes infrastructure specific details of a runnable service
// release, so that a container image can be mapped to its source code. This
// artifact is further used to properly update the service release within the
// underlying infrastructure definition, e.g. CloudFormation templates. The
// artifacts identity is equal to the cache index of the service release, which
// contains, amongst other information, the service ID of this artifact as
// Github or Docker repository, e.g. kayron or specta.
type Artifact struct {
	// Container holds information for whether the container image exists
	// according to the desired Git ref. The container image should exist if the
	// respective container image was pushed to the configured registry.
	Container container.Struct

	// Reference holds the current and desired Git ref for either a release tag or
	// a commit sha. This runnable codepoint may have to be derived from the
	// deployment strategy of any given service, if that deployment stratgy
	// defines a branch name, because branch names are not runnable without docker
	// image tag.
	Reference reference.Struct
}
