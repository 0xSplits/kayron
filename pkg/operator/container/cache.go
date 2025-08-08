package container

func (c *Container) cache(ima []image) {
	for _, x := range c.cac.Services() {
		var tag string
		{
			tag = curTag(ima, x.Release.Docker.String())
		}

		if tag == "" {
			continue
		}

		c.log.Log(
			"level", "debug",
			"message", "caching current state",
			"docker", x.Release.Docker.String(),
			"current", tag,
		)

		{
			x.Artifact.Scheduler.Current = tag
		}

		{
			c.cac.Update(x)
		}
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
