package cache

func (c *Cache) Releases() []Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	var lis []Object

	{
		lis = append(lis, c.inf...)
		lis = append(lis, c.ser...)
	}

	return lis
}
