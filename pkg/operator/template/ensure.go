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
		roo, err = t.rooSta()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Once we have found the root stack for this environment, we can set the
	// current state for this particular release artifact.

	var ver string
	{
		ver = temVer(roo.Parameters)
	}

	t.log.Log(
		"level", "debug",
		"message", "caching current state",
		"docker", inf.Release.Docker.String(),
		"current", ver,
	)

	{
		inf.Artifact.Scheduler.Current = ver
	}

	{
		t.cac.Update(inf)
	}

	return nil
}

func temVer(par []types.Parameter) string {
	for _, x := range par {
		if aws.ToString(x.ParameterKey) == "KayronTemplateVersion" {
			return aws.ToString(x.ParameterValue)
		}
	}

	return ""
}
