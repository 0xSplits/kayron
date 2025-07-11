package version

import (
	"github.com/spf13/cobra"
)

const (
	use = "version"
	sho = "Print the version information for this command line tool."
	lon = "Print the version information for this command line tool."
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
