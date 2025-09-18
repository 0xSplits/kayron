package provider

import (
	"github.com/0xSplits/kayron/pkg/constant"
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
	if s != constant.Cloudformation {
		return tracer.Mask(providerNameError)
	}

	return nil
}
