package infrastructure

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

func (i *Infrastructure) Ensure() error {
	var inf cache.Object
	{
		inf = i.cac.Infrastructure()
	}

	{
		i.log.Log(
			"level", "debug",
			"message", "resolved ref for github repository",
			"environment", i.env,
			"repository", fmt.Sprintf("https://github.com/%s/%s", i.own, inf.Release.Github.String()),
			"ref", inf.Artifact.Reference.Desired,
		)
	}

	// Setup the Github file system implementation so that we can walk the
	// contents of the remote infrastructure repository.

	var gfs *roghfs.Roghfs
	{
		gfs = roghfs.New(roghfs.Config{
			Bas: afero.NewMemMapFs(),
			Git: i.git,
			Own: i.own,
			Rep: inf.Release.Github.String(),
			Ref: inf.Artifact.Reference.Desired,
		})
	}

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
			byt, err = afero.ReadFile(gfs, pat)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		{
			err = i.putObj(pat, byt)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		return nil
	}

	{
		err := afero.Walk(gfs, Directory, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Once all infrastructure templates have been uploaded to S3, we can set the
	// condition success for this particular release artifact.

	{
		inf.Artifact.Condition.Success = true
	}

	{
		i.cac.Update(inf)
	}

	return nil
}
