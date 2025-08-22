package operator

import (
	"github.com/0xSplits/workit/handler"
)

func (o *Operator) Chain() [][]handler.Ensure {
	return [][]handler.Ensure{
		// Lookup all the release settings for the configured releases
		// repository. This first step initializes the release and artifact
		// caches and ensures that no cached state carries over from previous
		// reconciliation loops.
		{o.release},

		// Run the next steps in parallel in order to find the current and
		// desired state of the release artifacts that we are tasked to
		// managed.
		//
		//     1. Lookup the ECS container of every service release regardless
		//        of its deployment strategy. This populates the CURRENT state
		//        of the service artifacts.
		//
		//     2. Lookup the Github reference of every release definition that
		//        defines a branch deployment strategy. This populates the
		//        DESIRED state of all artifact references.
		//
		//     3. Lookup the CloudFormation template of the infrastructure
		//        release regardless of its deployment strategy. This
		//        populates the CURRENT state of the infrastructure artifact.
		//
		{o.container, o.reference, o.template},

		// Check whether those ECR image tags exist that are specified in the
		// desired state of any given service release. We only need to do this
		// for the service releases that have to get updated, which is why
		// this step must run after fetching the current and desired state of
		// our service releases.
		{o.registry},

		// Check whether we have any valid state drift amongst our cached
		// service releases. If we cannot detect any drift, then we do not
		// have to do any more work during this particular reconciliation
		// loop. This policy implementation is a control flow primitive with
		// the ability to cancel the reconciliation loop.
		{o.policy},

		// Once the current and desired states of the runnable service
		// releases are known to have drifted apart, we can fetch the current
		// version of our cloudformation templates from the configured
		// infrastructure repository. We only need to do this if there is at
		// least one service release that has to get updated.
		{o.infrastructure},

		// TODO document
		{o.cloudFormation},
	}
}
