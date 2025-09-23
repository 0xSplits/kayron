package preview

import (
	"context"
	"sort"
	"strings"

	"github.com/0xSplits/kayron/pkg/hash"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/branch"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/preview"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/tracer"
)

var (
	// filter is a collection of branch name prefixes that we want to ignore when
	// expanding a service release into preview releases. E.g. we do not want to
	// deploy preview releases for dependabot branches.
	filter = []string{
		"dependabot/",
	}
)

func (p *Preview) Expand(rel release.Struct) (release.Slice, error) {
	var err error

	var opt *github.PullRequestListOptions
	{
		opt = &github.PullRequestListOptions{
			State: "open",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}
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
	// Before sorting, filter pull requests by branch names that we definitely
	// want to consider for preview releases. E.g. drop all dependabot branches
	// before sorting.

	var fil []*github.PullRequest

	for _, x := range pul {
		if !hasPre(x.GetHead().GetRef(), filter) {
			fil = append(fil, x)
		}
	}

	// Sort our filtered pull requests from oldest to newest, so that new pull
	// requests do not change the order of pull request specific preview
	// deployments. This is relevant because the order of preview releases created
	// below will define the priority settings of the ALB's listener rules.

	sort.Slice(fil, func(i, j int) bool {
		return fil[i].GetCreatedAt().Before(fil[j].GetCreatedAt().Time)
	})

	var lis release.Slice
	for _, x := range fil {
		var pre release.Struct
		{
			pre = rel
		}

		var bra string
		var pul int
		{
			bra = x.GetHead().GetRef()
			pul = x.GetNumber()
		}

		{
			pre.Deploy = deploy.Struct{
				Branch:  branch.String(bra),
				Preview: preview.Bool(true),
			}
		}

		// Make sure to inject the preview deployment hash into the release labels.
		// This is used to identify the correct current state of deployed container
		// image tags, as well as rendering the correct CloudFormation templates.

		{
			pre.Labels.Hash = hash.New(bra)
			pre.Labels.Pull = pul
		}

		{
			lis = append(lis, pre)
		}
	}

	return lis
}

func hasPre(str string, pre []string) bool {
	for _, x := range pre {
		if strings.HasPrefix(str, x) {
			return true
		}
	}

	return false
}
