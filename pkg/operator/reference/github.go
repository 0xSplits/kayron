package reference

import (
	"context"
	"fmt"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/webhook"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/tracer"
)

func (r *Reference) desRef(rel release.Struct) (string, error) {
	// Resolve the commit sha if the branch deployment strategy is selected. We
	// try to fetch data on demand from Github's branch API first, and then try to
	// compare the fetched result against our internal webhook cache. If a newer
	// commit hash is present in our internal webhook cache, then we prefer that
	// one over the outdated/lagging version that Github's heavily cached API
	// endpoints provide. Note that branches may be referenced in releases while
	// they are not yet tracked, or not tracked anymore inside of Github. This may
	// happen predominantly during testing when preparing or finishing releases
	// and their dependencies.

	if !rel.Deploy.Branch.Empty() {
		pll, err := r.pllCom(rel.Github.String(), rel.Deploy.Branch.String())
		if err != nil {
			return "", tracer.Mask(err)
		}

		fmt.Printf("pll %#v\n", pll)

		psh, err := r.pshCom(rel.Github.String(), rel.Deploy.Branch.String(), pll)
		if err != nil {
			return "", tracer.Mask(err)
		}

		fmt.Printf("psh %#v\n", psh)

		return psh.Hash, nil
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

// pllCom tries to find the latest commit of a repository branch by pulling data
// from the Github API.
func (r *Reference) pllCom(rep string, ref string) (webhook.Commit, error) {
	var err error

	var bra *github.Branch
	var res *github.Response
	{
		bra, res, err = r.git.Repositories.GetBranch(context.Background(), r.own, rep, ref, 3)
		if isNotFound(res) {
			return webhook.Commit{}, nil
		} else if err != nil {
			return webhook.Commit{}, tracer.Mask(err,
				tracer.Context{Key: "owner", Value: r.own},
				tracer.Context{Key: "repository", Value: rep},
				tracer.Context{Key: "branch", Value: ref},
			)
		}
	}

	var com webhook.Commit
	{
		com = webhook.Commit{
			Hash: bra.GetCommit().GetSHA(),
			Time: bra.GetCommit().GetCommit().GetCommitter().GetDate().Time,
		}
	}

	return com, nil
}

// pshCom tries to find the latest commit of a repository branch by using the
// pushed data from a Github app webhook.
func (r *Reference) pshCom(rep string, ref string, pll webhook.Commit) (webhook.Commit, error) {
	var key webhook.Key
	{
		key = webhook.Key{
			Org: r.own,
			Rep: rep,
			Bra: ref,
		}
	}

	var com webhook.Commit
	{
		com = r.whk.Latest(key, pll)
	}

	if com.Empty() {
		r.log.Log(
			"level", "warning",
			"message", "git ref unresolvable",
			"reason", "branch not found",
			"suggestion", "this issue might be caused by a user error or eventual consistency of the underlying backend",
			"owner", r.own,
			"repository", rep,
			"branch", ref,
		)
	}

	return com, nil
}
