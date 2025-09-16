package cache

func (c *Cache) Update(obj Object) {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	if obj.kin == Infrastructure {
		c.inf[obj.ind].Artifact = c.inf[obj.ind].Artifact.Merge(obj.Artifact)
	}

	if obj.kin == Preview {
		c.pre[obj.ind].Artifact = c.pre[obj.ind].Artifact.Merge(obj.Artifact)
	}

	if obj.kin == Service {
		c.ser[obj.ind].Artifact = c.ser[obj.ind].Artifact.Merge(obj.Artifact)
	}
}
