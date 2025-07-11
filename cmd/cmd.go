package cmd

import (
	"github.com/0xSplits/kayron/cmd/daemon"
	"github.com/0xSplits/kayron/cmd/deploy"
	"github.com/0xSplits/kayron/cmd/lint"
	"github.com/0xSplits/kayron/cmd/version"
	"github.com/spf13/cobra"
)

var (
	use = "kayron"
	sho = "Golang based operator microservice."
	lon = "Golang based operator microservice."
)

func New() *cobra.Command {
	var cmd *cobra.Command
	{
		cmd = &cobra.Command{
			Use:   use,
			Short: sho,
			Long:  lon,
			Run:   (&run{}).run,
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
			// We slience errors because we do not want to see spf13/cobra printing.
			// The errors returned by the commands will be propagated to the main.go
			// anyway, where we have custom error printing for the command line
			// tool.
			SilenceErrors: true,
			SilenceUsage:  true,
		}
	}

	{
		cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	}

	{
		cmd.AddCommand(daemon.New())
		cmd.AddCommand(deploy.New())
		cmd.AddCommand(lint.New())
		cmd.AddCommand(version.New())
	}

	return cmd
}
