package loader

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/0xSplits/kayron/pkg/schema/specification"
	"github.com/0xSplits/kayron/pkg/schema/specification/labels"
	"github.com/0xSplits/kayron/pkg/schema/specification/service"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy"
	"github.com/spf13/afero"
)

// Test_Loader ensures that schemas can be loaded properly, according to the
// underlying file system.
func Test_Loader(t *testing.T) {
	testCases := []struct {
		roo string
		sch specification.Schemas
		err error
	}{
		// Case 000, no config
		{
			roo: ".",
			sch: nil,
			err: nil,
		},
		// Case 001, correct folder
		{
			roo: "./environment",
			sch: specification.Schemas{
				{
					Labels: labels.Labels{
						Environment: "production",
						Testing:     false,
					},
					Service: service.Services{
						{
							Docker: "kayron",
							Github: "kayron",
							Deploy: deploy.Deploy{
								Release: "v1.8.2",
							},
						},
					},
				},
				{
					Labels: labels.Labels{
						Environment: "staging",
						Testing:     false,
					},
					Service: service.Services{
						{
							Docker: "specta",
							Github: "specta",
							Deploy: deploy.Deploy{
								Suspend: true,
							},
						},
					},
				},
				{
					Labels: labels.Labels{
						Environment: "foobar",
						Testing:     true,
					},
					Service: service.Services{
						{
							Docker: "splits",
							Github: "server",
						},
					},
				},
				{
					Labels: labels.Labels{
						Environment: "testing",
						Testing:     true,
					},
					Service: service.Services{
						{
							Docker: "server",
							Github: "server",
							Deploy: deploy.Deploy{
								Branch: "feature",
							},
						},
					},
				},
			},
			err: nil,
		},
		// Case 002, folder name is "foo"
		{
			roo: "./environment",
			sch: nil,
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

			var sch specification.Schemas
			{
				sch, err = Loader(fil, roo)
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected %#v got %#v", tc.err, err)
				}
			}

			if !reflect.DeepEqual(sch, tc.sch) {
				t.Fatalf("expected %#v got %#v", tc.sch, sch)
			}
		})
	}
}
