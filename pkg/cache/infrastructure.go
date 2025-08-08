package cache

func (c *Cache) Infrastructure() Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	return c.inf[0] // Cache.Create guarantees 1 infrastructure release
}
