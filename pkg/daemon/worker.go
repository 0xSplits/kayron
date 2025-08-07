package daemon

import (
	"time"

	"github.com/0xSplits/kayron/pkg/context"
	"github.com/0xSplits/kayron/pkg/operator"
	"github.com/0xSplits/kayron/pkg/operator/policy"
	"github.com/0xSplits/workit/handler"
	"github.com/0xSplits/workit/worker"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/choreo/sequence"
)

func (d *Daemon) Worker() *worker.Worker {
	var cfg aws.Config
	{
		cfg = musAws()
	}

	var ctx *context.Context
	{
		ctx = context.New(context.Config{
			Log: d.log,
		})
	}

	var ope *operator.Operator
	{
		ope = operator.New(operator.Config{
			Aws: cfg,
			Ctx: ctx,
			Env: d.env,
			Log: d.log,
			Met: d.met,
		})
	}

	return worker.New(worker.Config{
		Env: d.env.Environment,
		Fil: policy.IsCancel,
		Han: []handler.Interface{
			handler.New(handler.Config{
				Coo: 10 * time.Second,
				Ens: sequence.Wrap(
					// Lookup all the release settings for the configured releases
					// repository. This first step initializes the release and artifact
					// caches and ensures that no cached state carries over from previous
					// reconciliation loops.
					ope.Release.Ensure,

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
					parallel.Wrap(ope.Container.Ensure, ope.Reference.Ensure, ope.Template.Ensure),

					// Check whether those ECR image tags exist that are specified in the
					// desired state of any given service release. We only need to do this
					// for the service releases that have to get updated, which is why
					// this step must run after fetching the current and desired state of
					// our service releases.
					ope.Registry.Ensure,

					// Check whether we have any valid state drift amongst our cached
					// service releases. If we cannot detect any drift, then we do not
					// have to do any more work during this particular reconciliation
					// loop. This policy implementation is a control flow primitive with
					// the ability to cancel the reconciliation loop.
					ope.Policy.Ensure,

					// Once the current and desired states of the runnable service
					// releases are known to have drifted apart, we can fetch the current
					// version of our cloudformation templates from the configured
					// infrastructure repository. We only need to do this if there is at
					// least one service release that has to get updated.
					ope.Infrastructure.Ensure,

					// TODO add business logic and document
					ope.CloudFormation.Ensure,
				),
			}),
		},
		Log: d.log,
		Met: d.met,
	})
}
