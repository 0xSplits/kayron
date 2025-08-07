package context

import (
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
)

func (c *Context) Create(rel release.Slice) {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	for _, x := range rel {
		c.log.Log(
			"level", "debug",
			"message", "caching release index",
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

		if x.Provider.String() == "cloudformation" {
			{
				obj.Kind = Infrastructure
			}

			{
				c.inf = append(c.inf, obj)
			}
		} else {
			{
				obj.Kind = Service
			}

			{
				c.ser = append(c.ser, obj)
			}
		}
	}
}
