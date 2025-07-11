package daemon

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xSplits/kayron/pkg/daemon"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/spf13/cobra"
)

type run struct {
	flag *flag
}

func (r *run) runE(cmd *cobra.Command, arg []string) error {
	var env envvar.Env
	{
		env = envvar.Load(r.flag.Env)
	}

	// --------------------------------------------------------------------- //

	var dae daemon.Interface
	{
		dae = daemon.New(daemon.Config{
			Env: env,
		})
	}

	if env.RunServer {
		go dae.Server().Daemon()
	}

	if env.RunWorker {
		go dae.Worker().Daemon()
	}

	// --------------------------------------------------------------------- //

	var sig chan os.Signal
	{
		sig = make(chan os.Signal, 2)
	}

	{
		defer close(sig)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	}

	{
		<-sig
	}

	select {
	case <-time.After(10 * time.Second):
		// One SIGTERM gives the daemon some time to tear down gracefully.
	case <-sig:
		// Two SIGTERMs stop the daemon immediately.
	}

	return nil
}
