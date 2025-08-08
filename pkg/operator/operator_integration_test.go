//go:build integration

package operator

import (
	"testing"

	"github.com/0xSplits/kayron/pkg/context"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/operator/policy"
	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/0xSplits/otelgo/recorder"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xh3b4sd/logger"
	"go.opentelemetry.io/otel/metric"
)

// Test_Operator_Integration runs the entire operator chain against all network
// dependencies, as if the operator was deploying for real. The required Github
// auth token to run this integration test needs at least public repo
// permissions. Also, a set of AWS credentials is required.
//
//	KAYRON_GITHUB_TOKEN=todo go test -tags=integration ./pkg/operator -v -race -run Test_Operator_Integration
func Test_Operator_Integration(t *testing.T) {
	var env envvar.Env
	{
		env = envvar.Env{
			Environment:   "testing",
			GithubToken:   envvar.MustGithub(),
			LogLevel:      "debug",
			ReleaseSource: "https://github.com/0xSplits/releases",
		}
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

	var ctx *context.Context
	{
		ctx = context.New(context.Config{
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

	var ope *Operator
	{
		ope = New(Config{
			Aws: cfg,
			Ctx: ctx,
			Dry: true, // dry run, read only
			Env: env,
			Log: log,
			Met: met,
		})
	}

	var fnc func() error
	{
		fnc = ope.Chain()
	}

	{
		err := fnc()
		if policy.IsCancel(err) {
			// fall through
		} else if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}
}
