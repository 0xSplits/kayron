package artifact

import (
	"github.com/0xSplits/kayron/pkg/release/artifact/container"
	"github.com/0xSplits/kayron/pkg/release/artifact/identity"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
)

// Artifact describes infrastructure specific details of a runnable service
// release, so that a container image can be mapped to its source code.
type Artifact struct {
	// Container contains the container image hash, if any. This runnable image
	// hash should exist if a container image was built and pushed for any given
	// Git reference, regardless its environment. If Oci is empty, then either the
	// service release is misconfigured, or the container image was not pushed.
	Container container.Struct

	// Identity contains the service ID of this artifact, e.g. kayron or specta
	// etc. This identifier is used to properly update the service release within
	// the underlying infrastructure, e.g. CloudFormation templates.
	Identity identity.String

	// Reference contains the Git reference of either a release tag or a commit
	// sha. This runnable codepoint has to be derived from the deployment strategy
	// of any given service, because deployment startegies may define branch
	// names, and branch names are not runnable.
	Reference reference.Struct
}
