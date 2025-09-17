package policy

import (
	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/cancel"
	"github.com/xh3b4sd/tracer"
)

func (p *Policy) Ensure() error {
	err := p.ensure(p.cac.Releases())
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

// ensure is just a private proxy for the public worker handler interface that
// is easier to test, because this method accepts the already parsed cache
// objects.
func (p *Policy) ensure(rel []cache.Object) error {
	// Verify that we have a valid artifact cache by iterating over all releases
	// and checking their cashed state before doing anything else. We must not
	// continue this reconciliation loop if there is any empty or invalid state,
	// because the side effects of proceeding using such a broken state could
	// potentially be dangerous. Note that verification is skipped for those
	// releases that are suspended.

	for _, x := range rel {
		if !bool(x.Release.Deploy.Suspend) && x.Artifact.Empty() {
			p.log.Log(
				"level", "warning",
				"message", "cancelling reconciliation loop",
				"reason", "invalid artifact cache",
				"docker", x.Release.Docker.String(),
				"github", x.Release.Github.String(),
				"provider", x.Release.Provider.String(),
				"current", musStr(x.Artifact.Scheduler.Current),
				"desired", musStr(x.Artifact.Reference.Desired),
			)

			return tracer.Mask(cacheStateEmptyError)
		}
	}

	// As soon as we detect a single valid state drift, we can return early and
	// allow the operator chain to execute the rest of the business logic. Our
	// current policy requires the following conditions to be true for a valid
	// state drift.
	//
	//     1. the desired deployment must not be suspended
	//
	//     2. the current and desired state must not be equal
	//
	//     3. the desired state must not be empty
	//
	//     4. the container image for the desired state must be pushed
	//

	var drf bool

	for _, x := range rel {
		if !bool(x.Release.Deploy.Suspend) && x.Artifact.Drift() && x.Artifact.Valid() {
			{
				drf = true
			}

			p.log.Log(
				"level", "info",
				"message", "continuing reconciliation loop",
				"reason", "detected state drift",
				"release", x.Name(),
				"domain", x.Domain(p.env.Environment),
				"version", x.Artifact.Reference.Desired,
			)
		}
	}

	if drf {
		return nil
	}

	// At this point all service releases were found to be up to date this time
	// around. This means that we do not have to do any more work for this
	// particular reconciliation loop. And so we return the control flow error
	// Cancel.

	p.log.Log(
		"level", "info",
		"message", "cancelling reconciliation loop",
		"reason", "no state drift",
	)

	return tracer.Mask(cancel.Error)
}

func musStr(str string) string {
	if str == "" {
		return "''"
	}

	return str
}
