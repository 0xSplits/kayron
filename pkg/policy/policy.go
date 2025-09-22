// Package policy inspects and verifies any potential state drift in order to
// either allow the reconciliation loop to continue, or cancel it. Only valid
// state drifts can be applied within the underlying infrastructure.
package policy

import (
	"fmt"
	"sync"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Aws aws.Config
	Cac *cache.Cache
	Env envvar.Env
	Log logger.Interface
}

type Policy struct {
	cac *cache.Cache
	cfc *cloudformation.Client
	env envvar.Env
	log logger.Interface
	mut sync.Mutex
	sta *types.Stack
}

func New(c Config) *Policy {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
	if c.Cac == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Cac must not be empty", c)))
	}
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Policy{
		cac: c.Cac,
		cfc: cloudformation.NewFromConfig(c.Aws),
		env: c.Env,
		log: c.Log,
		mut: sync.Mutex{},
		sta: nil,
	}
}
