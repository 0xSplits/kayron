package reference

import (
	"context"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/xh3b4sd/tracer"
)

func (r *Reference) desRef(rel release.Struct) (string, error) {
	// Return the commit sha if the branch deployment strategy is selected.

	if !rel.Deploy.Branch.Empty() {
		bra, _, err := r.git.Repositories.GetBranch(context.Background(), r.own, rel.Github.String(), rel.Deploy.Branch.String(), 3)
		if err != nil {
			return "", tracer.Mask(err)
		}

		return bra.GetCommit().GetSHA(), nil
	}

	// Return the configured release tag if the pinned release deployment strategy
	// is selected.

	if !rel.Deploy.Release.Empty() {
		return rel.Deploy.Release.String(), nil
	}

	// Fall through for e.g. suspended service deployments.
	//
	//     !rel.Deploy.Suspend.Empty()
	//

	return "", nil
}
