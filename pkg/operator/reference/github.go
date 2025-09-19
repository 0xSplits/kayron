// TODO rename file to git.go if this works
package reference

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/xh3b4sd/tracer"
)

func (r *Reference) desRef(rel release.Struct) (string, error) {
	// Return the commit sha if the branch deployment strategy is selected. Note
	// that branches may be referenced in releases while they are not yet tracked,
	// or not tracked anymore inside of Github. This may happen predominantly during
	// testing when preparing or finishing releases and their dependencies.

	if !rel.Deploy.Branch.Empty() {
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

func (r *Reference) comSha(rep string, bra string) (string, error) {
	var err error

	var rem *git.Remote
	{
		rem = git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
			URLs: []string{
				fmt.Sprintf("https://github.com/%s/%s.git", r.own, rep),
			},
		})
	}

	var opt *git.ListOptions
	{
		opt = &git.ListOptions{
			Auth: &http.BasicAuth{
				Username: "x-access-token",
				Password: r.env.GithubToken,
			},
		}
	}

	var ref []*plumbing.Reference
	{
		ref, err = rem.List(opt)
		if err != nil {
			return "", tracer.Mask(err,
				tracer.Context{Key: "owner", Value: r.own},
				tracer.Context{Key: "repository", Value: rep},
				tracer.Context{Key: "branch", Value: bra},
			)
		}
	}

	var nam plumbing.ReferenceName
	{
		nam = plumbing.NewBranchReferenceName(bra)
	}

	for _, x := range ref {
		if x.Name() == nam {
			return x.Hash().String(), nil
		}
	}

	r.log.Log(
		"level", "warning",
		"message", "git ref unresolvable",
		"reason", "branch not found",
		"suggestion", "this issue might be caused by a user error or eventual consistency of the underlying backend",
		"owner", r.own,
		"repository", rep,
		"branch", bra,
	)

	return "", nil
}
