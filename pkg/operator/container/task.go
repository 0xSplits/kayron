package container

import (
	"context"
	"slices"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

type task struct {
	// arn is the filtered task definition ARN that any given ECS service is
	// running right now.
	arn string
	// pre is the "preview" resource tag attached to any given ECS service, if
	// any, e.g. 1D0FD508.
	pre string
	// ser is the "service" resource tag attached to any given ECS service, e.g.
	// alloy or specta.
	ser string
}

// task resolves the task definition ARNs for the ECS services that are tagged
// using the "service" resource tag as defined by the provided list of service
// details.
func (c *Container) task(det []detail) ([]task, error) {
	var tas []task
	{
		tas = make([]task, len(det))
	}

	fnc := func(i int, d detail) error {
		var err error

		var inp *ecs.DescribeServicesInput
		{
			inp = &ecs.DescribeServicesInput{
				Cluster:  aws.String(d.clu),
				Services: []string{d.arn},
				Include:  []types.ServiceField{types.ServiceFieldTags},
			}
		}

		var out *ecs.DescribeServicesOutput
		{
			out, err = c.ecs.DescribeServices(context.Background(), inp)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		// Note that there should only ever be a single service in the response that
		// we iterate over below.

		if len(out.Services) != 1 {
			return tracer.Mask(invalidEcsServiceError)
		}

		for _, x := range out.Services {
			// For our reconciliation of the current state here to make sense we have
			// to ensure that the services being managed are actually controlled by
			// ECS. If we do not check this we may end up getting the wrong kind of
			// deployment information that we rely on to be provided in a certain way
			// below.

			if x.DeploymentController.Type != types.DeploymentControllerTypeEcs {
				c.log.Log(
					"level", "warning",
					"message", "skipping reconciliation for ECS service",
					"reason", "service does not use ECS controller",
					"cluster", *x.ClusterArn,
					"service", *x.ServiceArn,
				)

				{
					continue
				}
			}

			// There might be inactive or draining services with our desired service
			// labels in case we updated CloudFormation stacks multiple times during
			// with preview deployments during testing. We only want to consider the
			// current state of those stacks that are still active, because the
			// inactive versions have most likely been deleted already.

			if aws.ToString(x.Status) != "ACTIVE" {
				c.log.Log(
					"level", "debug",
					"message", "skipping reconciliation for ECS service",
					"reason", "ECS service is inactive",
					"cluster", *x.ClusterArn,
					"service", *x.ServiceArn,
				)

				{
					continue
				}
			}

			// Try to find the task definition that is successfully deployed. Here we
			// prevent picking those new task definitions prematurely that are still
			// transitioning during deployments.

			var arn string
			{
				arn = tasArn(x.Deployments)
			}

			if arn == "" {
				c.log.Log(
					"level", "debug",
					"message", "skipping reconciliation for ECS service",
					"reason", "ECS service has no completed task definition",
					"cluster", *x.ClusterArn,
					"service", *x.ServiceArn,
				)

				{
					continue
				}
			}

			// Try to identify the service and task definition based on their resource
			// tags.  If a service is not labelled with our desired 'service' tag,
			// then we cannot work with it properly moving forward.

			var pre string
			var ser string
			{
				pre = serTag(x.Tags, "preview")
				ser = serTag(x.Tags, "service")
			}

			if ser == "" {
				c.log.Log(
					"level", "warning",
					"message", "skipping reconciliation for ECS service",
					"reason", "ECS service has no 'service' tag",
					"cluster", *x.ClusterArn,
					"service", *x.ServiceArn,
				)

				{
					continue
				}
			}

			// This access pattern is concurrency safe because this callback only ever
			// executes for its own index.

			tas[i] = task{
				arn: arn,
				pre: pre,
				ser: ser,
			}
		}

		return nil
	}

	{
		err := parallel.Slice(det, fnc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	// Below we want to filter for the Docker image tags of those ECS services
	// that we have defined in the configured release source. In other words, if
	// there is no service release for e.g. specta, then we are not interested in
	// the Docker image tag of the specta containers in ECS. And so we can remove
	// them from our list and only fetch the image tags of the relevant service
	// releases.

	var ser []string
	for _, x := range c.cac.Services() {
		ser = append(ser, x.Release.Docker.String())
	}

	// Apply the service release filter according to the resource tag in AWS. Note
	// that the existance check below using slices is faster for small data sets
	// under roughly 50 items as compared to map checks that allocate.

	var fil []task
	for _, x := range tas {
		if slices.Contains(ser, x.ser) {
			fil = append(fil, x)
		}
	}

	return fil, nil
}

func serTag(tag []types.Tag, key string) string {
	for _, x := range tag {
		if aws.ToString(x.Key) == key {
			return aws.ToString(x.Value)
		}
	}

	return ""
}

func tasArn(dep []types.Deployment) string {
	// Make sure we can select the newest deployment first, so sort by time in
	// descending direction.

	sort.Slice(dep, func(i, j int) bool {
		return aws.ToTime(dep[i].CreatedAt).After(aws.ToTime(dep[j].CreatedAt))
	})

	// Now pick the task definition ARN of the first completed deployment.

	for _, y := range dep {
		if y.RolloutState == types.DeploymentRolloutStateCompleted {
			return aws.ToString(y.TaskDefinition)
		}
	}

	return ""
}
