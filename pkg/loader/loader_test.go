package loader

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/0xSplits/kayron/pkg/schema"
	"github.com/0xSplits/kayron/pkg/schema/service"
	"github.com/0xSplits/kayron/pkg/schema/service/deploy"
	"github.com/spf13/afero"
)

// Test_Loader ensures that schemas can be loaded properly, according to the
// underlying file system.
func Test_Loader(t *testing.T) {
	testCases := []struct {
		roo string
		sch schema.Schemas
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
			sch: schema.Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
							GitHub: "kayron",
							Deploy: deploy.Deploy{
								Release: "v1.8.2",
							},
						},
					},
				},
				{
					Service: service.Services{
						{
							Docker: "specta",
							GitHub: "specta",
							Deploy: deploy.Deploy{
								Suspend: true,
							},
						},
					},
				},
				{
					Service: service.Services{
						{
							Docker: "server",
							GitHub: "server",
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
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var fil afero.Fs
			{
				fil = afero.NewReadOnlyFs(afero.NewOsFs())
			}

			var roo string
			{
				roo = filepath.Join(".testdata", fmt.Sprintf("case.%03d", i), tc.roo)
			}

			var sch schema.Schemas
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
