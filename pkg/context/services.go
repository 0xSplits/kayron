package context

func (c *Context) Services() []Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	return c.ser
}
