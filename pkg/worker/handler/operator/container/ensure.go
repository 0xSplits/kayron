package container

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/xh3b4sd/tracer"
)

func (c *Container) Ensure() error {
	var err error

	var det []detail
	{
		det, err = c.detail()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	var tas []task
	{
		tas, err = c.task(det)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	var ima []image
	{
		ima, err = c.image(tas)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for i := range c.ser.Length() {
		var s service.Service
		{
			s, _ = c.ser.Search(i)
		}

		var tag string
		{
			tag = curTag(ima, s.Docker.String())
		}

		if tag != "" {
			fmt.Printf("%#v %#v\n", artifact.ReferenceCurrent(i), tag) // TODO emit proper logs
			c.art.Update(artifact.ReferenceCurrent(i), tag)
		}
	}

	return nil
}

func curTag(ima []image, ser string) string {
	for _, x := range ima {
		if x.ser == ser {
			return x.tag
		}
	}

	return ""
}
