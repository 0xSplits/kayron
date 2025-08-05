package container

import (
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
)

func (c *Container) cache(ima []image) {
	for i := range c.ser.Length() {
		var s service.Service
		{
			s, _ = c.ser.Search(i)
		}

		var tag string
		{
			tag = curTag(ima, s.Docker.String())
		}

		if tag == "" {
			continue
		}

		var key string
		{
			key = artifact.ContainerCurrent(i)
		}

		{
			c.art.Update(key, tag)
		}

		c.log.Log(
			"level", "debug",
			"message", "cached current state",
			"docker", s.Docker.String(),
			"github", s.Github.String(),
			"artifact", key,
			"current", tag,
		)
	}
}

func curTag(ima []image, ser string) string {
	for _, x := range ima {
		if x.ser == ser {
			return x.tag
		}
	}

	return ""
}
