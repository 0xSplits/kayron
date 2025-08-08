package cache

func (c *Cache) Infrastructure() Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	return c.inf[0]
}
