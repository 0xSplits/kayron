package registry

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/recorder"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Env envvar.Env
	Log logger.Interface

	Cou map[string]recorder.Interface
	Gau map[string]recorder.Interface
	His map[string]recorder.Interface
}

type Registry struct {
	env envvar.Env
	log logger.Interface

	cou map[string]recorder.Interface
	gau map[string]recorder.Interface
	his map[string]recorder.Interface
}

// New creates the fully configured registry whitelisting the set of counter
// metrics that are allowed to be tracked via the Kayron RPCs.
func New(c Config) *Registry {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	if c.Cou == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Cou must not be empty", c)))
	}
	if c.Gau == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Gau must not be empty", c)))
	}
	if c.His == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.His must not be empty", c)))
	}

	return &Registry{
		env: c.Env,
		log: c.Log,

		cou: c.Cou,
		gau: c.Gau,
		his: c.His,
	}
}
