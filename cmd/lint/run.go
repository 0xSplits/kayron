package lint

import (
	"github.com/0xSplits/kayron/pkg/constant"
	"github.com/0xSplits/kayron/pkg/release/loader"
	"github.com/0xSplits/kayron/pkg/release/schema"
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

	var sch schema.Schema
	{
		sch, err = loader.Loader(sys, r.flag.Pat, constant.Infrastructure, constant.Service)
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
