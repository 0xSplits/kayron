package specification

import (
	"github.com/0xSplits/kayron/pkg/generic"
	"github.com/xh3b4sd/tracer"
)

type Schemas []Schema

func (s Schemas) Empty() bool {
	return len(s) == 0
}

func (s Schemas) Verify() error {
	// Ensure that we use unique values for the infrastructure shorthands across
	// multiple environments. E.g. we must not deploy "prod" twice with different
	// configurations.
	{
		lis := generic.Duplicate(s.colSho())
		if len(lis) != 0 {
			return tracer.Mask(infrastructureShorthandError, tracer.Context{Key: "shorthand", Value: lis})
		}
	}

	if s.Empty() {
		return tracer.Mask(schemaEmptyError)
	}

	for _, x := range s {
		err := x.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (s Schemas) colSho() []string {
	var lis []string

	for _, x := range s {
		if x.Infrastructure.Shorthand != "" {
			lis = append(lis, x.Infrastructure.Shorthand)
		}
	}

	return lis
}
