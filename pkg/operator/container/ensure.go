package container

import (
	"github.com/xh3b4sd/tracer"
)

func (c *Container) Ensure() error {
	var err error

	// Fetch some details about the ECS services currently running in AWS.

	var det []detail
	{
		det, err = c.detail()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Fetch the ECS tasks associated to the running ECS services so that we can
	// inspect their container images.

	var tas []task
	{
		tas, err = c.task(det)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Find the Docker image tags based on the given ECS tasks.

	var ima []image
	{
		ima, err = c.image(tas)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Cache the current state of the configured service releases in terms of
	// their currently deployed version.

	{
		c.cache(ima)
	}

	return nil
}
