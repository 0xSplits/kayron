package template

import (
	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

func (t *Template) Ensure() error {
	var err error

	var inf cache.Object
	{
		inf = t.cac.Infrastructure()
	}

	var roo types.Stack
	{
		roo, err = t.sta.Search()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Once we have found the root stack for this environment, we can set the
	// current state and the condition success for this particular release
	// artifact.

	var ver string
	{
		ver = temVer(roo, inf)
	}

	t.log.Log(
		"level", "debug",
		"message", "caching current state",
		"github", inf.Release.Github.String(),
		"current", musStr(ver),
	)

	{
		inf.Artifact.Condition.Success = true
		inf.Artifact.Scheduler.Current = ver
	}

	{
		t.cac.Update(inf)
	}

	return nil
}

func musStr(str string) string {
	if str == "" {
		return "''"
	}

	return str
}

func temVer(roo types.Stack, inf cache.Object) string {
	for _, x := range roo.Parameters {
		if aws.ToString(x.ParameterKey) == inf.Parameter() {
			return aws.ToString(x.ParameterValue)
		}
	}

	return ""
}
