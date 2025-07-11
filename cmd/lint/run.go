package lint

import (
	"fmt"

	"github.com/spf13/cobra"
)

type run struct{}

func (r *run) run(cmd *cobra.Command, arg []string) {
	fmt.Printf("TODO\n")
}
