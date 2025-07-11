package daemon

import (
	"github.com/spf13/cobra"
)

const (
	use = "daemon"
	sho = "Execute Kayron's long running process for running the operator."
	lon = "Execute Kayron's long running process for running the operator."
)

func New() *cobra.Command {
	var flg *flag
	{
		flg = &flag{}
	}

	var cmd *cobra.Command
	{
		cmd = &cobra.Command{
			Use:   use,
			Short: sho,
			Long:  lon,
			RunE:  (&run{flag: flg}).runE,
		}
	}

	{
		flg.Init(cmd)
	}

	return cmd
}
