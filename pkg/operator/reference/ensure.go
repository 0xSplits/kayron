package reference

import (
	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

func (r *Reference) Ensure() error {
	var err error

	// Get the list of cached releases so that we can lookup their respective
	// artifact references concurrently, if necessary. This includes
	// infrastructure and service releases.

	var rel []cache.Object
	{
		rel = r.cac.Releases()
	}

	// Find the reference for every branch deployment strategy. The concurrently
	// executed function below prevents network calls for every release that does
	// not define a branch deployment strategy.

	fnc := func(i int, x cache.Object) error {
		ref, err := r.desRef(x.Release)
		if err != nil {
			return tracer.Mask(err)
		}

		if ref == "" {
			return nil
		}

		r.log.Log(
			"level", "debug",
			"message", "caching desired state",
			"github", x.Release.Github.String(),
			"desired", ref,
		)

		{
			x.Artifact.Reference.Desired = ref
		}

		{
			r.cac.Update(x)
		}

		return nil
	}

	{
		err = parallel.Slice(rel, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
