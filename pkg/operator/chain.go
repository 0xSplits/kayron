package operator

import (
	"github.com/0xSplits/workit/handler"
)

// Chain returns the directed acyclic graph of worker handlers executed by the
// sequential worker engine. All handlers of the same row are executed
// concurrently, while one row is executed sequentially, one after another.
func (o *Operator) Chain() [][]handler.Ensure {
	return [][]handler.Ensure{
		// Lookup all the release settings for the configured releases
		// repository. This first step initializes the release and artifact
		// caches and ensures that no cached state carries over from previous
		// reconciliation loops.
		{o.release},

		// Inject any potential preview deployments into our internal list of
		// release definitions so that we can render and expose any additional
		// development services during testing. Note that this operator function is
		// only active within the testing environment.
		{o.preview},

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

		// Run the next steps in parallel and split the execution graph based on the
		// different underlying responsibilities.
		//
		//     1. Once the current and desired states of the runnable service
		//        releases are known to have drifted apart, we can fetch the current
		//        version of our CloudFormation templates from the configured
		//        infrastructure repository, and upload all templates to S3. Note
		//        that this operator function is only active in case of valid state
		//        drift.
		//
		//     2. Emit log messages about the internal state of the system. This
		//        operator function runs on each and every reconciliation loop.
		//
		//     3. Manage any relevant deployment status for visibility purposes.
		//        E.g. create and update pull request comments about any preview
		//        deployment status.
		//
		{o.infrastructure, o.logging, o.status},

		// With the CloudFormation templates uploaded to S3, we can construct the
		// payload to update the CloudFormation stack that we are responsible for.
		// Optional parameters and tags will be considered for the UpdateStack
		// command. Once a new update cycle got initiated, the reconciliation loop
		// ends, and the operator should not reconcile the watched CloudFormation
		// stack again until the ongoing stack update ended up in some retryable
		// stack status. Note that this operator function is only active in case of
		// valid state drift.
		{o.cloudFormation},
	}
}
