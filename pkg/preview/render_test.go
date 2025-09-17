package preview

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/hash"
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/branch"
	"github.com/0xSplits/kayron/pkg/release/schema/release/docker"
	"github.com/0xSplits/kayron/pkg/release/schema/release/labels"
	"github.com/google/go-cmp/cmp"
)

func Test_Preview_Render(t *testing.T) {
	testCases := []struct {
		obj []cache.Object
	}{
		// Case 000
		{
			obj: []cache.Object{
				{
					Artifact: artifact.Struct{
						Reference: reference.Struct{
							Desired: "bc7891268e44f62e0aebbe339c0850b61d52c417",
						},
					},
					Release: release.Struct{
						Deploy: deploy.Struct{
							Branch: branch.String("fancy-feature-branch"),
						},
						Docker: docker.String("lite"),
						Labels: labels.Struct{
							Hash: hash.New("fancy-feature-branch"),
						},
					},
				},
				{
					Artifact: artifact.Struct{
						Reference: reference.Struct{
							Desired: "02b42b7ec63d4078767cb3b7cb0d34fde91b6237",
						},
					},
					Release: release.Struct{
						Deploy: deploy.Struct{
							Branch: branch.String("dependabot/another-one"),
						},
						Docker: docker.String("lite"),
						Labels: labels.Struct{
							Hash: hash.New("dependabot/another-one"),
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var err error

			var inp []byte
			{
				inp, err = os.ReadFile(fmt.Sprintf("./testdata/%03d/inp.yaml.golden", i))
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			var out []byte
			{
				out, err = os.ReadFile(fmt.Sprintf("./testdata/%03d/out.yaml.golden", i))
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			var pre *Preview
			{
				pre = New(Config{
					Env: envvar.Env{
						Environment:   "testing",
						GithubToken:   "foo",
						ReleaseSource: "https://github.com/0xSplits/releases",
					},
					Inp: inp,
				})
			}

			var res []byte
			{
				res, err = pre.Render(tc.obj)
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			if dif := cmp.Diff(bytes.TrimSpace(out), bytes.TrimSpace(res)); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
