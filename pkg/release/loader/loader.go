package loader

import (
	"io/fs"
	"path/filepath"

	"github.com/0xSplits/kayron/pkg/release/schema"
	specification "github.com/0xSplits/kayron/pkg/release/schema"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

// Loader collects all available service settings by walking the provided
// filesystem, starting at the given root directory. All .yaml files are
// inspected and their content is unmarshalled into slices of service.Service
// types. An error is returned if walking, reading or unmarshalling fails.
func Loader(sys afero.Fs, roo string) (schema.Schema, error) {
	var sch specification.Schema

	fnc := func(pat string, fil fs.FileInfo, err error) error {
		{
			if err != nil {
				return tracer.Mask(err)
			}
			if fil.IsDir() {
				return nil
			}
		}

		var ext string
		{
			ext = filepath.Ext(fil.Name())
			if ext != ".yaml" {
				return nil
			}
		}

		var byt []byte
		{
			byt, err = afero.ReadFile(sys, pat)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		var lis service.Slice
		{
			err = yaml.Unmarshal(byt, &lis)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		for i := range lis {
			lis[i].Labels.Source = pat
		}

		{
			sch.Service = append(sch.Service, lis...)
		}

		return nil
	}

	{
		err := afero.Walk(sys, filepath.Clean(roo), fnc)
		if err != nil {
			return schema.Schema{}, tracer.Mask(err)
		}
	}

	return sch, nil
}
