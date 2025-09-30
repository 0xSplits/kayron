package logging

import (
	"github.com/0xSplits/kayron/pkg/cache"
)

func (l *Logging) Ensure() error {
	var can bool
	{
		can = l.pol.Cancel()
	}

	if can {
		l.log.Log(
			"level", "info",
			"message", "deployment in progress",
		)

		return nil
	}

	// Find all release artifacts with state drift, whether their underlying
	// conditions are successful or not.

	var dft []cache.Object
	{
		dft = l.pol.Drift(false)
	}

	// If we do not have any drifted release artifacts, then this means that all
	// service releases were found to be up to date this time around.

	if len(dft) == 0 {
		l.log.Log(
			"level", "info",
			"message", "no state drift",
		)
	}

	// We may find release artifacts to show state drift, while some of their
	// underlying conditions may not yet be fulfilled. Note that the only
	// condition we have to wait for at the moment is the container image to be
	// pushed to the underlying container registry.

	for _, x := range dft {
		var msg string
		var rsn string

		if x.Drift(true) {
			msg = "verifying state drift"
			rsn = "all conditions successfull"
		} else if x.Drift(false) {
			msg = "detected state drift"
			rsn = "waiting for container image"
		}

		l.log.Log(
			"level", "info",
			"message", msg,
			"reason", rsn,
			"release", x.Name(),
			"preview", x.Release.Labels.Hash.Upper(),
			"domain", x.Domain(l.env.Environment),
			"version", x.Artifact.Reference.Desired,
		)
	}

	return nil
}
