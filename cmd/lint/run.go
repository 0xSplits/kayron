package lint

import (
	"github.com/0xSplits/kayron/pkg/schema/loader"
	"github.com/0xSplits/kayron/pkg/schema/specification"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type run struct {
	flag *flag
}

func (r *run) runE(cmd *cobra.Command, arg []string) error {
	var err error

	var sys afero.Fs
	{
		sys = afero.NewReadOnlyFs(afero.NewOsFs())
	}

	var sch specification.Schemas
	{
		sch, err = loader.Loader(sys, r.flag.Pat)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		err = sch.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
