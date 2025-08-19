package reference

import (
	"context"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/xh3b4sd/tracer"
)

func (r *Reference) desRef(rel release.Struct) (string, error) {
	// Return the commit sha if the branch deployment strategy is selected.

	if !rel.Deploy.Branch.Empty() {
		sha, err := r.comSha(rel)
		if err != nil {
			return "", tracer.Mask(err)
		}

		return sha, nil
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

func (r *Reference) comSha(rel release.Struct) (string, error) {
	bra, res, err := r.git.Repositories.GetBranch(context.Background(), r.own, rel.Github.String(), rel.Deploy.Branch.String(), 3)
	if isNotFound(res) {
		r.log.Log(
			"level", "warning",
			"message", "git ref unresolvable",
			"reason", "branch not found",
			"suggestion", "this issue might be caused by a user error or eventual consistency of the underlying backend",
			"owner", r.own,
			"repository", rel.Github.String(),
			"branch", rel.Deploy.Branch.String(),
		)

		return "", nil
	} else if err != nil {
		return "", tracer.Mask(err,
			tracer.Context{Key: "owner", Value: r.own},
			tracer.Context{Key: "repository", Value: rel.Github.String()},
			tracer.Context{Key: "branch", Value: rel.Deploy.Branch.String()},
		)
	}

	return bra.GetCommit().GetSHA(), nil
}
