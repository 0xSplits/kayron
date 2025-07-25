package daemon

import (
	"context"
	"os"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator"
	"github.com/0xSplits/kayron/pkg/worker/handler/releases"
	"github.com/0xSplits/workit/handler"
	"github.com/0xSplits/workit/worker"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Worker() *worker.Worker {
	var cfg aws.Config
	{
		cfg = musAws()
	}

	// The cache implementations used here are the dumb pipes connecting our
	// worker handlers by providing critical information required to make the
	// entire operator chain work as expected. All worker handlers are executed
	// iteratively all the time, but the real work may only be done if the cached
	// information provided is sufficient. E.g. an empty cache result may cause
	// some worker handlers to skip all work by returning early.

	var rel cache.Interface[int, service.Service]
	{
		rel = cache.New[int, service.Service]()
	}

	return worker.New(worker.Config{
		Env: d.env.Environment,
		Han: []handler.Interface{
			releases.New(releases.Config{Env: d.env, Log: d.log, Rel: rel}),
			operator.New(operator.Config{Aws: cfg, Env: d.env, Log: d.log, Met: d.met}),
		},
		Log: d.log,
		Met: d.met,
	})
}

func musAws() aws.Config {
	reg := os.Getenv("AWS_REGION")
	if reg == "" {
		reg = "us-west-2"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(reg))
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}

	return cfg
}
