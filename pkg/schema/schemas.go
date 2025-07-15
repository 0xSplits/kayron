package schema

import "github.com/xh3b4sd/tracer"

type Schemas []Schema

func (s Schemas) Empty() bool {
	return len(s) == 0
}

func (s Schemas) Verify() error {
	// TODO validate unique values across multiple environments, e.g. must not use
	// prod twice as env shorthand

	for _, x := range s {
		err := x.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
