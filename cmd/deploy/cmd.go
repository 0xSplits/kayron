package deploy

import (
	"github.com/spf13/cobra"
)

const (
	use = "deploy"
	sho = "Manually trigger a CloudFormation stack update."
	lon = "Manually trigger a CloudFormation stack update."
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
