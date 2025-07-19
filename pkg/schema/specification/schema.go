// Package specification defines the in-memory representation of Kayron’s
// change-management definitions.
package specification

import (
	"github.com/0xSplits/kayron/pkg/schema/specification/infrastructure"
	"github.com/0xSplits/kayron/pkg/schema/specification/labels"
	"github.com/0xSplits/kayron/pkg/schema/specification/service"
	"github.com/xh3b4sd/tracer"
)

// Schema is the root of our deployment configuration.
type Schema struct {
	Infrastructure infrastructure.Infrastructure `yaml:"infrastructure,omitempty"`
	Labels         labels.Labels                 `yaml:"-"`
	Service        service.Services              `yaml:"service,omitempty"`
}

func (s Schema) Empty() bool {
	return s.Infrastructure.Empty() && s.Service.Empty()
}

func (s Schema) Verify() error {
	if s.Empty() {
		return tracer.Mask(schemaEmptyError, tracer.Context{Key: "file", Value: s.Labels.Source})
	}

	{
		err := s.Infrastructure.Verify()
		if err != nil {
			return tracer.Mask(err, tracer.Context{Key: "file", Value: s.Labels.Source})
		}
	}

	{
		err := s.Service.Verify()
		if err != nil {
			return tracer.Mask(err, tracer.Context{Key: "file", Value: s.Labels.Source})
		}
	}

	return nil
}
