package reference

import (
	"context"

	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

func (r *Reference) Ensure() error {
	var err error

	// Get the list of cached service releases so that we can lookup their
	// respective artifact references concurrently, if necessary.

	var ser []service.Service
	for i := range r.ser.Length() {
		var s service.Service
		{
			s, _ = r.ser.Search(i)
		}

		{
			ser = append(ser, s)
		}
	}

	// Find the reference for every branch deployment strategy. The concurrently
	// executed function below prevents network calls for every service that does
	// not define a branch deployment strategy. Note that we can update the
	// indexed cache keys concurrently, because we are only ever updating cache
	// leafs, which is to say non-nested data structures.

	fnc := func(i int, x service.Service) error {
		ref, err := r.desRef(x)
		if err != nil {
			return tracer.Mask(err)
		}

		if ref == "" {
			return nil
		}

		var key string
		{
			key = artifact.ReferenceDesired(i)
		}

		{
			r.art.Update(key, ref)
		}

		r.log.Log(
			"level", "debug",
			"message", "cached desired state",
			"docker", x.Docker.String(),
			"github", x.Github.String(),
			"artifact", key,
			"desired", ref,
		)

		return nil
	}

	{
		err = parallel.Slice(ser, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (r *Reference) desRef(ser service.Service) (string, error) {
	// Return the commit sha if the branch deployment strategy is selected.

	if !ser.Deploy.Branch.Empty() {
		bra, _, err := r.git.Repositories.GetBranch(context.Background(), r.own, ser.Github.String(), ser.Deploy.Branch.String(), 3)
		if err != nil {
			return "", tracer.Mask(err)
		}

		return bra.GetCommit().GetSHA(), nil
	}

	// Return the configured release tag if the pinned release deployment strategy
	// is selected.

	if !ser.Deploy.Release.Empty() {
		return ser.Deploy.Release.String(), nil
	}

	// Fall through for e.g. suspended service deployments.
	//
	//     !ser.Deploy.Suspend.Empty()
	//

	return "", nil
}
