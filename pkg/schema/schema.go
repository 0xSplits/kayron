// Package schema defines the in-memory representation of Kayron’s
// change-management definitions.
package schema

import (
	"github.com/0xSplits/kayron/pkg/schema/infrastructure"
	"github.com/0xSplits/kayron/pkg/schema/service"
	"github.com/xh3b4sd/tracer"
)

// Schema is the root of our deployment configuration.
type Schema struct {
	Infrastructure infrastructure.Infrastructure `yaml:"infrastructure,omitempty"`
	Service        service.Services              `yaml:"service,omitempty"`
}

func (s Schema) Empty() bool {
	return s.Infrastructure.Empty() && s.Service.Empty()
}

func (s Schema) Verify() error {
	{
		err := s.Infrastructure.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		err := s.Service.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
