package version

import (
	"fmt"
	"io"
	"os"

	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type run struct{}

func (r *run) run(_ *cobra.Command, _ []string) {
	mustPrint(os.Stdout, "Git Sha       %s\n", runtime.Sha())
	mustPrint(os.Stdout, "Git Tag       %s\n", runtime.Tag())
	mustPrint(os.Stdout, "Repository    %s\n", runtime.Src())
	mustPrint(os.Stdout, "Go Arch       %s\n", runtime.Arc())
	mustPrint(os.Stdout, "Go OS         %s\n", runtime.Gos())
	mustPrint(os.Stdout, "Go Version    %s\n", runtime.Ver())
}

func mustPrint(w io.Writer, f string, a ...any) {
	_, err := fmt.Fprintf(w, f, a...)
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}
}
