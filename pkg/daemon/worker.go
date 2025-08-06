package daemon

import (
	"context"
	"os"
	"time"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/operator/cloudformation"
	"github.com/0xSplits/kayron/pkg/operator/container"
	"github.com/0xSplits/kayron/pkg/operator/infrastructure"
	"github.com/0xSplits/kayron/pkg/operator/policy"
	"github.com/0xSplits/kayron/pkg/operator/reference"
	"github.com/0xSplits/kayron/pkg/operator/registry"
	"github.com/0xSplits/kayron/pkg/operator/release"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/workit/handler"
	"github.com/0xSplits/workit/worker"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/choreo/sequence"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Worker() *worker.Worker {
	var cfg aws.Config
	{
		cfg = musAws()
	}

	// The cache implementations used here are the dumb pipes connecting our
	// business logic across different boundary conditions by providing critical
	// information required to make the entire operator chain work as expected.
	// The operator's worker handler is executed constantly in an interval, but
	// the real work may only be done if the cached information provided deeper
	// down the stack is sufficient. E.g. an empty cache result may cause some
	// business logic to be skipped temporarily. Note that the cache keys setup
	// here are cross references to the respective cache values within the various
	// cache implementations.

	var art cache.Interface[string, string]
	{
		art = cache.New[string, string]()
	}

	var ser cache.Interface[int, service.Service]
	{
		ser = cache.New[int, service.Service]()
	}

	var clo *cloudformation.CloudFormation
	var con *container.Container
	var inf *infrastructure.Infrastructure
	var pol *policy.Policy
	var ref *reference.Reference
	var rel *release.Release
	var reg *registry.Registry
	{
		clo = cloudformation.New(cloudformation.Config{Art: art, Aws: cfg, Env: d.env, Log: d.log, Met: d.met, Ser: ser})
		con = container.New(container.Config{Art: art, Aws: cfg, Env: d.env, Log: d.log, Ser: ser})
		inf = infrastructure.New(infrastructure.Config{Art: art, Aws: cfg, Env: d.env, Log: d.log, Ser: ser})
		pol = policy.New(policy.Config{Art: art, Log: d.log, Ser: ser})
		ref = reference.New(reference.Config{Art: art, Env: d.env, Log: d.log, Ser: ser})
		rel = release.New(release.Config{Art: art, Env: d.env, Log: d.log, Ser: ser})
		reg = registry.New(registry.Config{Art: art, Aws: cfg, Env: d.env, Log: d.log, Ser: ser})
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
					rel.Ensure,

					// Run the next steps in parallel in order to find the current and
					// desired state of the service releases that we are tasked to
					// managed.
					//
					//     1. Lookup the ECS container of every service release regardless
					//        of its deployment strategy. This populates the CURRENT state
					//        of the artifact reference.
					//
					//     2. Lookup the Github reference of every service release that
					//        defines a branch deployment strategy. This populates the
					//        DESIRED state of the artifact reference.
					//
					parallel.Wrap(con.Ensure, ref.Ensure),

					// Check whether those ECR image tags exist that are specified in the
					// desired state of any given service release. We only need to do this
					// for the service releases that have to get updated, which is why
					// this step must run after fetching the current and desired state of
					// our service releases.
					reg.Ensure,

					// Check whether we have any valid state drift amongst our cached
					// service releases. If we cannot detect any drift, then we do not
					// have to do any more work during this particular reconciliation
					// loop. This policy implementation is a control flow primitive with
					// the ability to cancel the reconciliation loop.
					pol.Ensure,

					// Once the current and desired states of the runnable service
					// releases are known to have drifted apart, we can fetch the current
					// version of our cloudformation templates from the configured
					// infrastructure repository. We only need to do this if there is at
					// least one service release that has to get updated.
					inf.Ensure,

					// TODO add business logic and document
					clo.Ensure,
				),
			}),
		},
		Log: d.log,
		Met: d.met,
	})
}

func musAws() aws.Config {
	reg := os.Getenv("AWS_REGION")
	if reg == "" {
		reg = "us-west-2"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(reg))
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}

	return cfg
}
