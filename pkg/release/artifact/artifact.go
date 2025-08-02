// Package artifact provides cache key formatting for infrastructure specific
// details of a runnable service release, so that a container image can be
// mapped to its source code.  This artifact is further used to properly update
// the service release within the underlying infrastructure definition, e.g.
// CloudFormation templates. The artifacts identity is equal to the cache index
// of the service release, which contains, amongst other information, the
// service ID of this artifact as Github or Docker repository, e.g. kayron or
// specta.
package artifact

import (
	"fmt"
)

// The Container namespace keys hold information for whether the container image
// exists according to the desired Git ref. The container image should exist if
// the respective container image was pushed to the configured registry.

// ContainerExists returns the indexed key for the desired artifact reference
// for which there exists a Docker image tag inside of the configured Docker
// registry.
func ContainerExists(i int) string {
	return fmt.Sprintf("[%d].Container.Exists", i)
}

// The Reference namespace keys hold the current and desired Git ref for either
// a release tag or a commit sha. This runnable codepoint may have to be derived
// from the deployment strategy of any given service, if that deployment stratgy
// defines a branch name, because branch names are not runnable without docker
// image tag.

// ReferenceCurrent returns the indexed key for the image tag currently deployed
// within the underlying infrastructure provider, e.g. CloudFormation.
func ReferenceCurrent(i int) string {
	return fmt.Sprintf("[%d].Reference.Current", i)
}

// ReferenceDesired returns the indexed key for either the release tag or commit
// sha for this artifact, depending on the specified deployment startegy.
func ReferenceDesired(i int) string {
	return fmt.Sprintf("[%d].Reference.Desired", i)
}
