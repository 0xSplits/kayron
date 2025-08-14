package canceler

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/stack"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Aws aws.Config
	Env envvar.Env
	Sta stack.Interface
}

type Canceler struct {
	cfc *cloudformation.Client
	env envvar.Env
	sta stack.Interface
}

func New(c Config) *Canceler {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Sta == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Sta must not be empty", c)))
	}

	return &Canceler{
		cfc: cloudformation.NewFromConfig(c.Aws),
		env: c.Env,
		sta: c.Sta,
	}
}
