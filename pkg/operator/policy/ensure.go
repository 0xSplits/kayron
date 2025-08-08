package policy

import (
	"github.com/xh3b4sd/tracer"
)

func (p *Policy) Ensure() error {
	// Verify that we have a valid artifact cache by iterating over all releases
	// and checking their cashed state before doing anything else. We must not
	// continue this reconciliation loop if there is any empty or invalid state,
	// because the side effects of proceeding using such a broken state could
	// potentially be dangerous.

	for _, x := range p.ctx.Releases() {
		if x.Artifact.Empty() {
			p.log.Log(
				"level", "warning",
				"message", "cancelling reconciliation loop",
				"reason", "invalid artifact cache",
				"docker", x.Release.Docker.String(),
				"github", x.Release.Github.String(),
				"provider", x.Release.Provider.String(),
				"current", x.Artifact.Scheduler.Current,
				"desired", x.Artifact.Reference.Desired,
			)

			return tracer.Mask(cacheStateEmptyError)
		}
	}

	// As soon as we detect a single valid state drift, we can return early and
	// allow the operator chain to execute the rest of the business logic. Our
	// current policy requires the following conditions to be true for a valid
	// state drift.
	//
	//     1. the current and desired state must not be equal
	//
	//     2. the desired state must not be empty
	//
	//     3. the container image for the desired state must be pushed
	//

	for _, x := range p.ctx.Releases() {
		if x.Artifact.Drift() && x.Artifact.Valid() {
			p.log.Log(
				"level", "info",
				"message", "continuing reconciliation loop",
				"reason", "detected state drift",
			)

			return nil
		}
	}

	// At this point all service releases were found to be up to date this time
	// around. This means that we do not have to do any more work for this
	// particular reconciliation loop. And so we return the control flow error
	// Cancel.

	p.log.Log(
		"level", "debug",
		"message", "cancelling reconciliation loop",
		"reason", "no state drift",
	)

	return tracer.Mask(cancelError)
}
