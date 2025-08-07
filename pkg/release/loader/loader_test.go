package loader

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/0xSplits/kayron/pkg/release/schema"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/docker"
	"github.com/0xSplits/kayron/pkg/release/schema/release/github"
	"github.com/0xSplits/kayron/pkg/release/schema/release/labels"
	"github.com/0xSplits/kayron/pkg/release/schema/release/provider"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

// Test_Loader ensures that release definitions can be loaded properly,
// according to the underlying file system.
func Test_Loader(t *testing.T) {
	testCases := []struct {
		roo string
		sch schema.Schema
		err error
	}{
		// Case 000, no config
		{
			roo: ".",
			sch: schema.Schema{},
			err: nil,
		},
		// Case 001, root folder
		{
			roo: ".",
			sch: schema.Schema{
				Release: release.Slice{
					{
						Github:   github.String("infrastructure"),
						Provider: provider.String("cloudformation"),
						Deploy: deploy.Struct{
							Release: "v1.13.0",
						},
						Labels: labels.Struct{
							Source: "testdata/case.001/infrastructure/infrastructure.yaml",
						},
					},
					{
						Docker: docker.String("splits"),
						Github: github.String("server"),
						Labels: labels.Struct{
							Source: "testdata/case.001/service/foobar.yaml",
						},
					},
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "testdata/case.001/service/kayron.yaml",
						},
					},
					{
						Docker: docker.String("server"),
						Github: github.String("server"),
						Deploy: deploy.Struct{
							Branch: "feature",
						},
						Labels: labels.Struct{
							Source: "testdata/case.001/service/server.yaml",
						},
					},
					{
						Docker: docker.String("specta"),
						Github: github.String("specta"),
						Deploy: deploy.Struct{
							Suspend: true,
						},
						Labels: labels.Struct{
							Source: "testdata/case.001/service/specta.yaml",
						},
					},
				},
			},
			err: nil,
		},
		// Case 002, folder name is "foo"
		{
			roo: ".",
			sch: schema.Schema{
				Release: release.Slice{
					{
						Docker: docker.String("splits"),
						Github: github.String("server"),
						Labels: labels.Struct{
							Source: "testdata/case.002/bar/foobar.yaml",
						},
					},
					{
						Docker: docker.String("server"),
						Github: github.String("server"),
						Deploy: deploy.Struct{
							Branch: "feature",
						},
						Labels: labels.Struct{
							Source: "testdata/case.002/bar/testing.yaml",
						},
					},
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "testdata/case.002/baz/kayron.yaml",
						},
					},
					{
						Docker: docker.String("specta"),
						Github: github.String("specta"),
						Deploy: deploy.Struct{
							Suspend: true,
						},
						Labels: labels.Struct{
							Source: "testdata/case.002/baz/specta.yaml",
						},
					},
					{
						Github:   github.String("infrastructure"),
						Provider: provider.String("cloudformation"),
						Deploy: deploy.Struct{
							Release: "v1.13.0",
						},
						Labels: labels.Struct{
							Source: "testdata/case.002/foo/infrastructure.yaml",
						},
					},
				},
			},
			err: nil,
		},
		// Case 003, folder "environment" does not exist
		{
			roo: "./environment",
			sch: schema.Schema{},
			err: fs.ErrNotExist,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var err error

			var fil afero.Fs
			{
				fil = afero.NewReadOnlyFs(afero.NewOsFs())
			}

			var roo string
			{
				roo = filepath.Join("testdata", fmt.Sprintf("case.%03d", i), tc.roo)
			}

			var sch schema.Schema
			{
				sch, err = Loader(fil, roo)
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected %#v got %#v", tc.err, err)
				}
			}

			if dif := cmp.Diff(tc.sch, sch); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
