package preview

import (
	"fmt"
	"testing"
	"time"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/branch"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/preview"
	"github.com/0xSplits/kayron/pkg/release/schema/release/docker"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v73/github"
)

func Test_Preview_Expand(t *testing.T) {
	testCases := []struct {
		rel release.Struct
		pul []*github.PullRequest
		exp release.Slice
	}{
		// Case 000
		{
			rel: release.Struct{
				Deploy: deploy.Struct{
					Branch:  branch.String("main"),
					Preview: preview.Bool(true),
				},
				Docker: docker.String("lite"),
			},
			pul: []*github.PullRequest{
				{CreatedAt: tesTim(3), Head: tesBra("b/3")},
				{CreatedAt: tesTim(5), Head: tesBra("b/5")},
			},
			exp: release.Slice{
				{
					Deploy: deploy.Struct{
						Branch:  branch.String("main"),
						Preview: preview.Bool(false),
					},
					Docker: docker.String("lite"),
				},
				{
					Deploy: deploy.Struct{
						Branch:  branch.String("b/3"),
						Preview: preview.Bool(true),
					},
					Docker: docker.String("lite"),
				},
				{
					Deploy: deploy.Struct{
						Branch:  branch.String("b/5"),
						Preview: preview.Bool(true),
					},
					Docker: docker.String("lite"),
				},
			},
		},
		// Case 001
		{
			rel: release.Struct{
				Deploy: deploy.Struct{
					Branch:  branch.String("main"),
					Preview: preview.Bool(true),
				},
				Docker: docker.String("lite"),
			},
			pul: []*github.PullRequest{
				{CreatedAt: tesTim(3), Head: tesBra("b/3")},
				{CreatedAt: tesTim(5), Head: tesBra("b/5")},
				{CreatedAt: tesTim(7), Head: tesBra("b/7")},
				{CreatedAt: tesTim(9), Head: tesBra("b/9")},
			},
			exp: release.Slice{
				{
					Deploy: deploy.Struct{
						Branch:  branch.String("main"),
						Preview: preview.Bool(false),
					},
					Docker: docker.String("lite"),
				},
				{
					Deploy: deploy.Struct{
						Branch:  branch.String("b/3"),
						Preview: preview.Bool(true),
					},
					Docker: docker.String("lite"),
				},
				{
					Deploy: deploy.Struct{
						Branch:  branch.String("b/5"),
						Preview: preview.Bool(true),
					},
					Docker: docker.String("lite"),
				},
				{
					Deploy: deploy.Struct{
						Branch:  branch.String("b/7"),
						Preview: preview.Bool(true),
					},
					Docker: docker.String("lite"),
				},
				{
					Deploy: deploy.Struct{
						Branch:  branch.String("b/9"),
						Preview: preview.Bool(true),
					},
					Docker: docker.String("lite"),
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			exp := expand(tc.rel, tc.pul)

			if dif := cmp.Diff(tc.exp, exp); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}

func tesBra(nam string) *github.PullRequestBranch {
	return &github.PullRequestBranch{Ref: github.Ptr(nam)}
}

func tesTim(sec int64) *github.Timestamp {
	return &github.Timestamp{Time: time.Unix(sec, 0)}
}
