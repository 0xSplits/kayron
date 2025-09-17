package cache

// Previews returns all cached service release artifacts that are defined as
// preview deployments.
func (c *Cache) Previews(doc string) []Object {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	var lis []Object

	for _, x := range c.ser {
		if bool(x.Release.Deploy.Preview) && x.Release.Docker.String() == doc {
			lis = append(lis, x)
		}
	}

	return lis
}
