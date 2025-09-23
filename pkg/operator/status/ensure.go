package status

import (
	"github.com/xh3b4sd/tracer"
)

func (s *Status) Ensure() error {
	// Manage status updates for preview deployments in Github pull requests.

	{
		err := s.preview()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
