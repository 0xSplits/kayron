package infrastructure

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/constant"
	"github.com/0xSplits/kayron/pkg/preview"
	"github.com/0xSplits/roghfs"
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
			"environment", i.env.Environment,
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

		// Skip everything that is not a YAML file.

		var ext string
		{
			ext = filepath.Ext(fil.Name())
			if ext != ".yaml" {
				return nil
			}
		}

		var nam string
		{
			nam = strings.TrimSuffix(fil.Name(), ext)
		}

		var byt []byte
		{
			byt, err = afero.ReadFile(gfs, pat)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		// Before uploading our templates to S3, we have to inject any preview
		// deployments configured for the service release that matches this
		// particular template by file name.
		//
		//     "splits-lite" == Release.Docker.String()
		//

		{
			byt, err = i.renPre(nam, byt)
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
		err := afero.Walk(gfs, constant.Cloudformation, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (i *Infrastructure) renPre(nam string, byt []byte) ([]byte, error) {
	var err error

	// If there are no preview deployments defined inside our release artifacts,
	// then we do not have to inject anything, but instead return the same
	// template bytes early that we just received as input.

	var rel []cache.Object
	{
		rel = i.cac.Previews(nam)
	}

	if len(rel) == 0 {
		return byt, nil
	}

	// At this point we have preview deployments defined by at least one release
	// artifact. So we create a preview renderer and extend the raw template bytes
	// that we received as input data above.

	var pre *preview.Preview
	{
		pre = preview.New(preview.Config{
			Env: i.env,
			Git: i.git,
			Inp: byt,
		})
	}

	{
		byt, err = pre.Render(rel)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return byt, nil
}
