package container

func (c *Container) cache(ima []image) {
	for _, x := range c.cac.Services() {
		var tag string
		{
			tag = curTag(ima, x.Release.Docker.String())
		}

		c.log.Log(
			"level", "debug",
			"message", "caching current state",
			"docker", x.Release.Docker.String(),
			"current", musStr(tag),
		)

		// It may happen that there is no tag for services that are deployed the
		// first time. In those cases we want to log the current state above, but we
		// do not have to perform any artifact update. So we skip the loop below.

		if tag == "" {
			continue
		}

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

func musStr(str string) string {
	if str == "" {
		return "''"
	}

	return str
}
