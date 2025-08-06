package policy

import (
	"strconv"

	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/xh3b4sd/tracer"
)

func (p *Policy) Ensure() error {
	for i := range p.ser.Length() {
		var err error

		// Iterate over all service releases, except those referencing our
		// cloudformation templates. This exception may be removed in a later
		// refactoring, once we have a better idea about the differentiation between
		// services and infrastructure on the configuration level.

		var ser service.Service
		{
			ser, _ = p.ser.Search(i)
		}

		// TODO the way we mingle the infrastructure settings with our service
		// versions does not allow us to deploy infrastructure version changes. We
		// are currently excluding the infrastructure repo from the state drift
		// checks.
		if ser.Provider == "cloudformation" {
			continue
		}

		var exi string
		var cur string
		var des string
		{
			exi, _ = p.art.Search(artifact.ContainerDesired(i))
			cur, _ = p.art.Search(artifact.ContainerCurrent(i))
			des, _ = p.art.Search(artifact.ReferenceDesired(i))
		}

		if exi == "" || cur == "" || des == "" {
			// This should never happen, but if it did, we must not continue this
			// conciliation loop, because the side effects of proceeding using invalid
			// state could potentially be dangerous.

			p.log.Log(
				"level", "warning",
				"message", "invalid artifact cache",
				"reason", "detected empty state",
				"docker", ser.Docker.String(),
				"github", ser.Github.String(),
				"exists", exi,
				"current", cur,
				"desired", des,
			)

			return tracer.Mask(cancelError)
		}

		var ima bool
		{
			ima, err = strconv.ParseBool(exi)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		if ima && cur != des {
			// As soon as we detect a single valid state drift, we can return early
			// and allow the worker engine to execute the rest of the business logic.
			// Our current policy requires the following conditions to be true for a
			// valid state drift.
			//
			//     1. the current and desired state must not be equal
			//
			//     2. either current and desired state must not be empty
			//
			//     3. the container image for the new desired state must be pushed
			//

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
