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
