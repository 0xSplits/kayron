package cache

func (c *Cache) Previews(doc string) []Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	var lis []Object

	for _, x := range c.pre {
		if x.Release.Docker.String() == doc {
			lis = append(lis, x)
		}
	}

	return lis
}
