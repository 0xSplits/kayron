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

	// Get the list of cached service releases and release artifacts, so that we
	// can lookup the artifact references concurrently, if necessary.

	var art []artifact.Artifact
	var ser []service.Service
	for i := range r.art.Length() {
		var a artifact.Artifact
		var s service.Service
		{
			a, _ = r.art.Search(i)
			s, _ = r.ser.Search(i)
		}

		{
			art = append(art, a)
			ser = append(ser, s)
		}
	}

	// Find the reference for every branch deployment strategy. The concurrently
	// executed function below prevents network calls for every service that does
	// not define a branch deployment strategy.

	fnc := func(i int, x service.Service) error {
		var err error

		art[i].Reference.Desired, err = r.desRef(x)
		if err != nil {
			return tracer.Mask(err)
		}

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
