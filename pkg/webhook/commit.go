package webhook

import (
	"time"

	"github.com/xh3b4sd/tracer"
)

type Commit struct {
	Hash string
	Time time.Time
}

func (c Commit) Empty() bool {
	return c.Hash == "" && c.Time.IsZero()
}

func (c Commit) Equals(x Commit) bool {
	return c.Hash == x.Hash && c.Time.Equal(x.Time)
}

func (c Commit) Verify() error {
	if c.Hash == "" {
		return tracer.Mask(commitHashEmptyError)
	}

	if c.Time.IsZero() {
		return tracer.Mask(commitTimeEmptyError)
	}

	return nil
}
