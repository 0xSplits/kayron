package preview

import (
	"context"
	"sort"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/branch"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/preview"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/tracer"
)

func (p *Preview) Expand(rel release.Struct) (release.Slice, error) {
	var err error

	opt := &github.PullRequestListOptions{
		State: "open",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var pul []*github.PullRequest
	{
		pul, _, err = p.git.PullRequests.List(context.Background(), p.own, rel.Github.String(), opt)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var lis release.Slice
	{
		lis = expand(rel, pul)
	}

	return lis, nil
}

func expand(rel release.Struct, pul []*github.PullRequest) release.Slice {
	// Sort pull requests from oldest to newest, so that new pull requests do not
	// change the order of pull request specific preview deployments. This is
	// relevant because the order of preview releases created below will define
	// the priority settings of the ALB's listener rules.

	sort.Slice(pul, func(i, j int) bool {
		return pul[i].GetCreatedAt().Before(pul[j].GetCreatedAt().Time)
	})

	// Mark the expanded service release as non-preview. The Deploy.Preview
	// flag acts as a signal to expand our release definitions internally.
	// Once expanded, we redefine the purpose of this preview flag to maintain
	// our understanding of how to deploy "real" service releases. In other
	// words, we turn one release into many, while muting the one that
	// instructed the many for the preview mechanism.

	{
		rel.Deploy.Preview = preview.Bool(false)
	}

	var lis release.Slice
	{
		lis = append(lis, rel)
	}

	for _, x := range pul {
		var pre release.Struct
		{
			pre = rel
		}

		{
			pre.Deploy = deploy.Struct{
				Branch:  branch.String(x.GetHead().GetRef()),
				Preview: preview.Bool(true),
			}
		}

		{
			lis = append(lis, pre)
		}
	}

	return lis
}
