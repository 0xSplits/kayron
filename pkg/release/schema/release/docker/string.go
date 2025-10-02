package docker

import (
	"fmt"

	"github.com/distribution/reference"
	"github.com/xh3b4sd/tracer"
)

type String string

func (s String) Empty() bool {
	return s == ""
}

func (s String) String() string {
	return string(s)
}

func (s String) Verify() error {
	_, err := reference.ParseAnyReference(fmt.Sprintf("registry/%s:v0.1.0", s))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
