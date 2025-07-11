package version

import (
	"fmt"
	"os"

	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/spf13/cobra"
)

type run struct{}

func (r *run) run(cmd *cobra.Command, arg []string) {
	fmt.Fprintf(os.Stdout, "Git Sha       %s\n", runtime.Sha())
	fmt.Fprintf(os.Stdout, "Git Tag       %s\n", runtime.Tag())
	fmt.Fprintf(os.Stdout, "Repository    %s\n", runtime.Src())
	fmt.Fprintf(os.Stdout, "Go Arch       %s\n", runtime.Arc())
	fmt.Fprintf(os.Stdout, "Go OS         %s\n", runtime.Gos())
	fmt.Fprintf(os.Stdout, "Go Version    %s\n", runtime.Ver())
}
