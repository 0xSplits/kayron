package cache

// Services returns all cached service release artifacts, including those of any
// preview deployments.
func (c *Cache) Services() []Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	return c.ser
}
