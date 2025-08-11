package cache

import (
	"github.com/0xSplits/kayron/pkg/constant"
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/xh3b4sd/tracer"
)

func (c *Cache) Create(rel release.Slice) error {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	for _, x := range rel {
		c.log.Log(
			"level", "debug",
			"message", "caching release artifact",
			"docker", x.Docker.String(),
			"github", x.Github.String(),
			"deploy", x.Deploy.String(),
			"provider", x.Provider.String(),
		)

		var obj Object
		{
			obj = Object{
				Artifact: artifact.Struct{},
				Release:  x,
			}
		}

		if x.Provider.String() == constant.Cloudformation {
			{
				obj.Artifact.Condition.Trigger = c.frc // forward the infrastructure trigger
			}

			{
				obj.ind = len(c.inf)
				obj.kin = Infrastructure
			}

			{
				c.inf = append(c.inf, obj)
			}
		} else {
			{
				obj.ind = len(c.ser)
				obj.kin = Service
			}

			{
				c.ser = append(c.ser, obj)
			}
		}
	}

	if len(c.inf) != 1 {
		return tracer.Mask(invalidInfrastructureError)
	}

	return nil
}
