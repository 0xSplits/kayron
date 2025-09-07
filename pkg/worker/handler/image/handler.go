package image

import (
	"fmt"
	"strings"
	"time"

	"github.com/0xSplits/kayron/pkg/awsauthn"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/xh3b4sd/choreo/jitter"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	// Keep is the minimum amount of the most recent images that we want to keep
	// in the underlying registry, per repository, at all times.
	Keep = 10

	// Drop is the maximum amount of older images that we want to cleanup at once,
	// while respecting the minimum amount of images to keep at all times.
	Drop = 2
)

type Config struct {
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Pre prefix
	Rep repository
}

// Handler implements handler.Cooler as a cleanup job for container images. A
// single handler is responsible for cleaning up a single specific combination
// of container repository and image tag prefixes, e.g. kayron:v*.
type Handler struct {
	ecr *ecr.Client
	env envvar.Env
	jit *jitter.Jitter[time.Duration]
	key *awsauthn.Keychain
	log logger.Interface
	pre []string
	rep repository
}

func New(c Config) *Handler {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Pre == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Pre must not be empty", c)))
	}
	if c.Rep == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Rep must not be empty", c)))
	}

	return &Handler{
		ecr: ecr.NewFromConfig(c.Aws),
		env: c.Env,
		jit: jitter.New[time.Duration](jitter.Config{Per: 0.01}),
		key: awsauthn.New(awsauthn.Config{Aws: c.Aws}),
		log: c.Log,
		pre: strings.Split(string(c.Pre), ","),
		rep: c.Rep,
	}
}
