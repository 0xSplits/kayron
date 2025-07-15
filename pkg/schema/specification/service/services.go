package service

import "github.com/xh3b4sd/tracer"

// Service describes a list of independently deployable units, defined by a
// Docker image and a GitHub repository, plus the deployment strategies used for
// their rollout.
type Services []Service

func (s Services) Empty() bool {
	return len(s) == 0
}

func (s Services) Verify() error {
	for _, x := range s {
		err := x.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
