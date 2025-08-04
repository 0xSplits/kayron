package operator

import (
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/policy"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Ensure() error {
	var err error

	// Lookup all the release settings for the configured releases repository.
	// This first step initializes the release and artifact caches.

	{
		err = h.rel.Ensure()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Run the next steps in parallel in order to find the current and desired
	// state of the service releases that we are tasked to managed.
	//
	//     1. Lookup the ECS container of every service release regardless of its
	//        deployment strategy. This populates the CURRENT state of the
	//        artifact reference.
	//
	//     2. Lookup the Github reference of every service release that defines a
	//        branch deployment strategy. This populates the DESIRED state of the
	//        artifact reference.
	//

	{
		err = parallel.Func(h.con.Ensure, h.ref.Ensure)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Check whether we have any drift amongst our cached service releases. If we
	// cannot detect any drift, then we do not have to do any more work during
	// this particular reconciliation loop. Note that this policy implementation
	// is a control flow primitive that should eventually be supported by our
	// worker engine library.

	{
		err = h.pol.Ensure()
		if policy.IsCancel(err) {
			return nil
		} else if err != nil {
			return tracer.Mask(err)
		}
	}

	// Once the current and desired states of the runnable service releases are
	// known, we can run the next steps in parallel too. Additionally, we only
	// need to do real work in those following steps, if we recognize any state
	// drift.
	//
	//     1. Check whether those ECR image tags exist that are specified in the
	//        desired state of any given service release. We only need to do this
	//        for the service releases that have to get updated.
	//
	//     2. Fetch the current version of our cloudformation templates from the
	//        configured infrastructure repository. We only need to do this if
	//        there is at least one service release that has to get updated.
	//

	{
		err = parallel.Func(h.reg.Ensure, h.inf.Ensure)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		err = h.clo.Ensure() // TODO add business logic
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
