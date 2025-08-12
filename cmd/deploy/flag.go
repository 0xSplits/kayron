package deploy

import (
	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type flag struct {
	Env string
	Frc bool
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Env, "env", ".env", "the environment file to load")
	cmd.Flags().BoolVar(&f.Frc, "force", false, "whether to force the stack update")
}

func (f *flag) Validate() error {
	if f.Env == "" {
		return tracer.Mask(runtime.ExecutionFailedError, tracer.Context{Key: "reason", Value: "--env must not be empty"})
	}

	return nil
}
