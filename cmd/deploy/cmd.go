package deploy

import (
	"github.com/spf13/cobra"
)

const (
	use = "deploy"
	sho = "Claim the given test environment and continuously deploy the given service branch."
	lon = "Claim the given test environment and continuously deploy the given service branch."
)

func New() *cobra.Command {
	var cmd *cobra.Command
	{
		cmd = &cobra.Command{
			Use:   use,
			Short: sho,
			Long:  lon,
			Run:   (&run{}).run,
		}
	}

	return cmd
}
