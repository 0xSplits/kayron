package status

import (
	"github.com/0xSplits/kayron/pkg/cache"
)

func (s *Status) Ensure() error {
	var can bool
	{
		can = s.pol.Cancel()
	}

	if can {
		s.log.Log(
			"level", "info",
			"message", "deployment in progress",
		)

		return nil
	}

	var dft []cache.Object
	{
		dft = s.pol.Drift()
	}

	// If we do not have any drifted release artifacts, then this means that all
	// service releases were found to be up to date this time around.

	if len(dft) == 0 {
		s.log.Log(
			"level", "info",
			"message", "no state drift",
		)
	}

	// If we do have drifted release artifacts, then we want to create an info log
	// message for each affected service release.

	for _, x := range dft {
		s.log.Log(
			"level", "info",
			"message", "detected state drift",
			"release", x.Name(),
			"preview", x.Release.Labels.Hash.Upper(),
			"domain", x.Domain(s.env.Environment),
			"version", x.Artifact.Reference.Desired,
		)
	}

	// TODO manage Github comments/notifications

	return nil
}
