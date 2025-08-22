package deploy

import (
	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/cancel"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/operator"
	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/0xSplits/kayron/pkg/stack"
	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/workit/registry"
	"github.com/0xSplits/workit/worker/sequence"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

type run struct {
	flag *flag
}

func (r *run) runE(cmd *cobra.Command, arg []string) error {
	var env envvar.Env
	{
		env = envvar.Load(r.flag.Env)
	}

	var log logger.Interface
	{
		log = logger.New(logger.Config{
			Filter: logger.NewLevelFilter(env.LogLevel),
			Format: logger.JSONIndenter,
		})
	}

	var cfg aws.Config
	{
		cfg = envvar.MustAws()
	}

	var cac *cache.Cache
	{
		cac = cache.New(cache.Config{
			Frc: r.flag.Frc,
			Log: log,
		})
	}

	var met metric.Meter
	{
		met = recorder.NewMeter(recorder.MeterConfig{
			Env: env.Environment,
			Sco: "kayron",
			Ver: runtime.Tag(),
		})
	}

	var sta stack.Interface
	{
		sta = stack.New(stack.Config{
			Aws: cfg,
			Env: env,
			Log: log,
		})
	}

	var ope *operator.Operator
	{
		ope = operator.New(operator.Config{
			Aws: cfg,
			Cac: cac,
			Env: env,
			Log: log,
			Met: met,
			Sta: sta,
		})
	}

	var reg *registry.Registry
	{
		reg = registry.New(registry.Config{
			Env: env.Environment,
			Fil: cancel.Is,
			Log: log,
			Met: met,
		})
	}

	var wor *sequence.Worker
	{
		wor = sequence.New(sequence.Config{
			Han: ope.Chain(),
			Log: log,
			Reg: reg,
		})
	}

	{
		err := wor.Ensure()
		if cancel.Is(err) {
			// fall through
		} else if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
