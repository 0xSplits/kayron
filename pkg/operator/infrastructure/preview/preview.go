package preview

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/scanner"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Inp []byte
}

type Preview struct {
	inp []byte
	sca *scanner.Scanner
}

func New(c Config) *Preview {
	if c.Inp == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Inp must not be empty", c)))
	}

	return &Preview{
		inp: c.Inp,
		sca: scanner.New(scanner.Config{
			Inp: c.Inp,
		}),
	}
}
