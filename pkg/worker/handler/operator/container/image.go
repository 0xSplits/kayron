package container

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/distribution/reference"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

type image struct {
	// ser is the "service" resource tag attached to any given ECS service, e.g.
	// alloy or specta.
	ser string
	// tag is the Docker image tag that the task definition of any given ECS
	// service is running right now.
	tag string
}

// image resolves the current Docker image tag for any given task definition.
func (c *Container) image(tas []task) ([]image, error) {
	var err error

	var ima []image
	{
		ima = make([]image, len(tas))
	}

	fnc := func(i int, t task) error {
		var inp *ecs.DescribeTaskDefinitionInput
		{
			inp = &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: aws.String(t.arn),
			}
		}

		var out *ecs.DescribeTaskDefinitionOutput
		{
			out, err = c.ecs.DescribeTaskDefinition(context.Background(), inp)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		for _, x := range out.TaskDefinition.ContainerDefinitions {
			tag, err := imaTag(*x.Image)
			if err != nil {
				return tracer.Mask(err)
			}

			ima[i] = image{
				ser: t.ser,
				tag: tag,
			}
		}

		return nil
	}

	{
		err = parallel.Slice(tas, fnc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return ima, nil
}

// imaTag returns whatever version string a Docker image reference carries.
// imaTag understands every legal reference form that Docker itself accepts
// because it delegates all parsing to the official distribution/reference
// package. Note that imaTag does not default to "latest", because we have to
// resolve the exact image tag as used in EC2.
func imaTag(str string) (string, error) {
	ref, err := reference.ParseAnyReference(str)
	if err != nil {
		return "", tracer.Mask(err)
	}

	if t, ok := ref.(reference.Tagged); ok {
		return t.Tag(), nil
	}

	if d, ok := ref.(reference.Digested); ok {
		return d.Digest().String(), nil
	}

	// It is important to not default to the "latest" tag. If we did, then we
	// would trigger deployments over and over again for any ECS service that is
	// not specifying any particular image tag.

	return "", nil
}
