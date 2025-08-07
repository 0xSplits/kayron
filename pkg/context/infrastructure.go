package context

func (c *Context) Infrastructure() Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	return c.inf[0]
}
