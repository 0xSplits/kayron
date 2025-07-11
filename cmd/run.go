package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type run struct{}

func (r *run) run(cmd *cobra.Command, arg []string) {
	err := cmd.Help()
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}
}
