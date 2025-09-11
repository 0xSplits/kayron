// Package scanner is a multi line block scanner for YAML input bytes.
package scanner

import (
	"fmt"

	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Inp []byte
}

type Scanner struct {
	inp []byte
}

func New(c Config) *Scanner {
	if c.Inp == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Inp must not be empty", c)))
	}

	return &Scanner{
		inp: c.Inp,
	}
}
