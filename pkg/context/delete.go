package context

func (c *Context) Delete() {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	c.log.Log(
		"level", "debug",
		"message", "resetting operator cache",
	)

	{
		c.inf = nil
		c.ser = nil
	}
}
