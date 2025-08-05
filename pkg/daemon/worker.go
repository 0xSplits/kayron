package daemon

import (
	"context"
	"os"

	"github.com/0xSplits/kayron/pkg/worker/handler/operator"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/policy"
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

	return worker.New(worker.Config{
		Env: d.env.Environment,
		Fil: policy.IsCancel,
		Han: []handler.Interface{
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
