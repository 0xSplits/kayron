package preview

import (
	"bytes"

	"github.com/goccy/go-yaml"
	"github.com/xh3b4sd/tracer"
)

func lisPri(lin []byte) (int, error) {
	var pri map[string]int

	{
		err := yaml.Unmarshal(bytes.TrimSpace(lin), &pri)
		if err != nil {
			return 0, tracer.Mask(err)
		}
	}

	return pri["Priority"], nil
}
