package daemon

import (
	"time"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/operator"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/0xSplits/kayron/pkg/webhook"
	"github.com/0xSplits/kayron/pkg/worker/handler/image"
	"github.com/0xSplits/workit/handler"
	"github.com/0xSplits/workit/registry"
	"github.com/0xSplits/workit/worker/combined"
	"github.com/0xSplits/workit/worker/parallel"
	"github.com/0xSplits/workit/worker/sequence"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func (d *Daemon) Worker() *combined.Worker {
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

	var pol *policy.Policy
	{
		pol = policy.New(policy.Config{
			Aws: cfg,
			Cac: cac,
			Env: d.env,
			Log: d.log,
		})
	}

	var whk *webhook.Webhook
	{
		whk = webhook.New(webhook.Config{
			Env: d.env,
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
			Pol: pol,
			Whk: whk,
		})
	}

	var reg *registry.Registry
	{
		reg = registry.New(registry.Config{
			Env: d.env.Environment,
			Log: d.log,
			Met: d.met,
		})
	}

	var par *parallel.Worker
	{
		par = parallel.New(parallel.Config{
			Han: []handler.Cooler{
				image.New(image.Config{Aws: cfg, Env: d.env, Log: d.log, Pre: image.Sha, Rep: image.Kayron}),
				image.New(image.Config{Aws: cfg, Env: d.env, Log: d.log, Pre: image.Tag, Rep: image.Kayron}),

				image.New(image.Config{Aws: cfg, Env: d.env, Log: d.log, Pre: image.Sha, Rep: image.Server}),
				image.New(image.Config{Aws: cfg, Env: d.env, Log: d.log, Pre: image.Tag, Rep: image.Server}),

				image.New(image.Config{Aws: cfg, Env: d.env, Log: d.log, Pre: image.Sha, Rep: image.Specta}),
				image.New(image.Config{Aws: cfg, Env: d.env, Log: d.log, Pre: image.Tag, Rep: image.Specta}),

				image.New(image.Config{Aws: cfg, Env: d.env, Log: d.log, Pre: image.Sha, Rep: image.SplitsLite}),
				image.New(image.Config{Aws: cfg, Env: d.env, Log: d.log, Pre: image.Tag, Rep: image.SplitsLite}),
			},
			Log: d.log,
			Reg: reg,
		})
	}

	var seq *sequence.Worker
	{
		seq = sequence.New(sequence.Config{
			Coo: 10 * time.Second,
			Han: ope.Chain(),
			Log: d.log,
			Reg: reg,
		})
	}

	var wor *combined.Worker
	{
		wor = combined.New(combined.Config{
			Par: par,
			Seq: seq,
		})
	}

	return wor
}
