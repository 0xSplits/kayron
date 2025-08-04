// Package artifact provides cache key formatting for infrastructure specific
// details of a runnable service release, so that a container image can be
// mapped to its source code. Those artifact information are used to update the
// services within the underlying infrastructure according to their release
// defintions. Every artifact's identity is equal to the cache index of the
// service release, which contains, amongst other information, the service ID as
// Github or Docker repository, e.g. kayron or specta.
package artifact

import (
	"fmt"
)

// The Container namespace keys hold information for whether the container image
// exists according to configured container registry and container runtime.

// ContainerCurrent returns the indexed key for the container image tag of any
// given service release DEPLOYED within the configured container RUNTIME.
func ContainerCurrent(i int) string {
	return fmt.Sprintf("[%d].Container.Current", i)
}

// ContainerDesired returns the indexed key for the container image tag of any
// given service release PUSHED to the configured container REGISTRY.
func ContainerDesired(i int) string {
	return fmt.Sprintf("[%d].Container.Desired", i)
}

// The Reference namespace keys holds information about the Git ref of any given
// service release, which is either a release tag or a commit sha. This runnable
// codepoint may have to be derived from the deployment strategy for any given
// service, if its deployment stratgy defines a branch name, because branch
// names are not runnable without container image tag.

// ReferenceDesired returns the indexed key for either the release tag or commit
// sha of this artifact, depending on the specified deployment strategy.
func ReferenceDesired(i int) string {
	return fmt.Sprintf("[%d].Reference.Desired", i)
}
