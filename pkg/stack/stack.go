package stack

import (
	"fmt"
	"sync"

	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
}

type Stack struct {
	cfc *cloudformation.Client
	sta *types.Stack
	env envvar.Env
	log logger.Interface
	mut sync.Mutex
}

func New(c Config) *Stack {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Stack{
		cfc: cloudformation.NewFromConfig(c.Aws),
		sta: nil,
		env: c.Env,
		log: c.Log,
		mut: sync.Mutex{},
	}
}
