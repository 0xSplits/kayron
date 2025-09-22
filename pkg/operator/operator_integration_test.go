//go:build integration

package operator

import (
	"testing"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/cancel"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/0xSplits/kayron/pkg/stack"
	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/workit/registry"
	"github.com/0xSplits/workit/worker/sequence"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xh3b4sd/logger"
	"go.opentelemetry.io/otel/metric"
)

// Test_Operator_Integration runs the entire operator chain against all network
// dependencies, as if the operator was deploying for real. The required Github
// auth token to run this integration test needs at least public repo
// permissions. Also, a set of standard AWS credentials is required with
// read-only access as described in the README.md.
//
//	KAYRON_GITHUB_TOKEN=todo go test -tags=integration ./pkg/operator -v -race -run Test_Operator_Integration
func Test_Operator_Integration(t *testing.T) {
	var env envvar.Env

	{
		env = envvar.Env{
			CloudformationStack: "server-test",
			Environment:         "testing",
			GithubToken:         envvar.MustGithub(),
			LogLevel:            "debug",
			ReleaseSource:       "https://github.com/0xSplits/releases",
			S3Bucket:            "splits-cf-templates",
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

	var cac *cache.Cache
	{
		cac = cache.New(cache.Config{
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

	var pol *policy.Policy
	{
		pol = policy.New(policy.Config{
			Cac: cac,
			Env: env,
			Log: log,
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

	var ope *Operator
	{
		ope = New(Config{
			Aws: cfg,
			Cac: cac,
			Dry: true, // dry run, read only
			Env: env,
			Log: log,
			Met: met,
			Pol: pol,
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
			t.Fatal("expected", nil, "got", err)
		}
	}
}
