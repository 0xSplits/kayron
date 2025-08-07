package context

func (c *Context) Update(obj Object) {
	{
		c.mut.Lock()
		defer c.mut.Unlock()
	}

	if obj.Kind == Infrastructure {
		for i, x := range c.inf {
			if x.Release.Github == obj.Release.Github {
				c.inf[i] = obj
				break
			}
		}
	}

	if obj.Kind == Service {
		for i, x := range c.ser {
			if x.Release.Docker == obj.Release.Docker {
				c.inf[i] = obj
				break
			}
		}
	}
}
