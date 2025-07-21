package lint

import (
	"github.com/spf13/cobra"
)

const (
	use = "lint"
	sho = "Validate the release configuration under the given path."
	lon = "Validate the release configuration under the given path."
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
