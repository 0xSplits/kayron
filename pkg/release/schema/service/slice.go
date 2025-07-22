package service

import "github.com/xh3b4sd/tracer"

// Slice describes a list of independently deployable services, defined by a
// Docker image and a GitHub repository, plus the deployment strategies used for
// their rollout.
type Slice []Service

func (s Slice) Empty() bool {
	return len(s) == 0
}

func (s Slice) Verify() error {
	{
		emp := s.Empty()
		if emp {
			return tracer.Mask(serviceDefinitionEmptyError)
		}
	}

	for _, x := range s {
		err := x.Verify()
		if err != nil {
			return tracer.Mask(err, tracer.Context{Key: "index", Value: x.Labels.Block})
		}
	}

	return nil
}
