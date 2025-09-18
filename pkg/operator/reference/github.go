package reference

import (
	"context"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/xh3b4sd/tracer"
)

func (r *Reference) desRef(rel release.Struct) (string, error) {
	// Return the commit sha if the branch deployment strategy is selected. Note
	// that branches may be referenced in releases while they are not yet tracked,
	// or not tracked anymore inside of Github. This may happen predominantly during
	// testing when preparing or finishing releases and their dependencies. Note
	// that we do not lookup branch references for preview deployments, if those
	// references got already filled in the release labels. E.g. we may have looked
	// up the latest commit sha for a preview deployment in an earlier stage of this
	// reconciliation loop already.

	if !rel.Deploy.Branch.Empty() {
		if bool(rel.Deploy.Preview) && rel.Labels.Head != "" {
			return rel.Labels.Head, nil
		}

		sha, err := r.comSha(rel.Github.String(), rel.Deploy.Branch.String())
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

func (r *Reference) comSha(rep string, ref string) (string, error) {
	bra, res, err := r.git.Repositories.GetBranch(context.Background(), r.own, rep, ref, 3)
	if isNotFound(res) {
		r.log.Log(
			"level", "warning",
			"message", "git ref unresolvable",
			"reason", "branch not found",
			"suggestion", "this issue might be caused by a user error or eventual consistency of the underlying backend",
			"owner", r.own,
			"repository", rep,
			"branch", ref,
		)

		return "", nil
	} else if err != nil {
		return "", tracer.Mask(err,
			tracer.Context{Key: "owner", Value: r.own},
			tracer.Context{Key: "repository", Value: rep},
			tracer.Context{Key: "branch", Value: ref},
		)
	}

	return bra.GetCommit().GetSHA(), nil
}
