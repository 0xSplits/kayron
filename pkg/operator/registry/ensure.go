package registry

import (
	"strconv"

	"github.com/0xSplits/kayron/pkg/context"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

func (r *Registry) Ensure() error {
	var ser []context.Object
	{
		ser = r.ctx.Services()
	}

	// Check whether the desired Docker image exists within the underlying
	// container registry, if the current and desired state differs.

	fnc := func(i int, x context.Object) error {
		var err error

		var cur string
		var des string
		{
			cur = x.Artifact.Scheduler.Current
			des = x.Artifact.Reference.Desired
		}

		// We do not have to do any work here if the currently deployed service
		// already matches the desired service release.

		if cur == des {
			return nil
		}

		var exi bool
		{
			exi, err = r.imaExi(x.Release.Docker.String(), des)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		{
			x.Artifact.Condition.Success = exi
		}

		{
			r.ctx.Update(x)
		}

		r.log.Log(
			"level", "debug",
			"message", "executed image check",
			"image", x.Release.Docker.String(),
			"tag", des,
			"exists", strconv.FormatBool(exi),
		)

		return nil
	}

	{
		err := parallel.Slice(ser, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
