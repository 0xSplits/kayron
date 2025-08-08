package cache

func (c *Cache) Services() []Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	return c.ser
}
