package container

import (
	"context"
	"slices"

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
			var tag string
			{
				tag = serTag(x.Tags)
			}

			if tag == "" {
				c.log.Log(
					"level", "warning",
					"message", "skipping instrumentation for ECS service",
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
				arn: *x.TaskDefinition,
				ser: tag,
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

func serTag(tag []types.Tag) string {
	for _, x := range tag {
		if *x.Key == "service" {
			return *x.Value
		}
	}

	return ""
}
