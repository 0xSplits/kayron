package daemon

import (
	"time"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/operator"
	"github.com/0xSplits/kayron/pkg/operator/policy"
	"github.com/0xSplits/workit/handler"
	"github.com/0xSplits/workit/worker"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func (d *Daemon) Worker() *worker.Worker {
	var cfg aws.Config
	{
		cfg = envvar.MustAws()
	}

	var cac *cache.Cache
	{
		cac = cache.New(cache.Config{
			Log: d.log,
		})
	}

	var ope *operator.Operator
	{
		ope = operator.New(operator.Config{
			Aws: cfg,
			Cac: cac,
			Env: d.env,
			Log: d.log,
			Met: d.met,
		})
	}

	return worker.New(worker.Config{
		Env: d.env.Environment,
		Fil: policy.IsCancel,
		Han: []handler.Interface{
			handler.New(handler.Config{
				Coo: 10 * time.Second,
				Ens: ope.Chain(),
			}),
		},
		Log: d.log,
		Met: d.met,
	})
}
