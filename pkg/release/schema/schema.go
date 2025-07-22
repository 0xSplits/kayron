// Package schema defines the in-memory representation of Kayron’s change
// management definitions.
package schema

import (
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/xh3b4sd/tracer"
)

// Schema is the root of our release configuration.
type Schema struct {
	Service service.Slice
}

func (s Schema) Empty() bool {
	return s.Service.Empty()
}

func (s Schema) Verify() error {
	err := s.Service.Verify()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
