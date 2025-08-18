package loader

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/0xSplits/kayron/pkg/release/schema"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

// Loader collects all available service releases by walking the provided
// filesystem, starting at the given root directory. All .yaml files are
// inspected and their content is unmarshalled into slices of release.Struct
// types. An error is returned if walking, reading or unmarshalling fails.
func Loader(sys afero.Fs, roo string, wht ...string) (schema.Schema, error) {
	var sch schema.Schema

	// Ensure that all provided file paths are stripped clean for consistency.

	{
		roo = filepath.Clean(roo)
	}

	for i := range wht {
		wht[i] = filepath.Clean(wht[i])
	}

	fnc := func(pat string, fil fs.FileInfo, err error) error {
		// Do some sanity checks for the current iteration. Any error stops the walk
		// immediately. If the given path is not whitelisted, then we skip its
		// associated folder. And then, any folder at all does not have to be
		// inspected, because we are looking for .yaml files.

		{
			if err != nil {
				return tracer.Mask(err)
			}
			if fil.IsDir() {
				if !whtPat(roo, wht, pat) {
					return fs.SkipDir // SkipDir must never be returned for files
				} else {
					return nil
				}
			}
		}

		var ext string
		{
			ext = filepath.Ext(fil.Name())
			if ext != ".yaml" {
				return nil
			}
		}

		// At this point we found a .yaml file and can read its byte content, which
		// we simply marshal into a list of release definitions.

		var byt []byte
		{
			byt, err = afero.ReadFile(sys, pat)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		var lis release.Slice
		{
			err = yaml.Unmarshal(byt, &lis)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		for i := range lis {
			lis[i].Labels.Block = i
			lis[i].Labels.Source = pat
		}

		{
			sch.Release = append(sch.Release, lis...)
		}

		return nil
	}

	{
		err := afero.Walk(sys, roo, fnc)
		if err != nil {
			return schema.Schema{}, tracer.Mask(err)
		}
	}

	return sch, nil
}

func whtPat(roo string, wht []string, pat string) bool {
	// The very first path that we are trafersing is the root path itself, which
	// should never be subject to skipping, because that would mean to not walk
	// the file system at all.

	if roo == pat {
		return true
	}

	// In case there is no whitelist specified, we return true so that we do not
	// skip walking this path.

	if len(wht) == 0 {
		return true
	}

	// At this point we have a whitelist injected and we have to respect it. All
	// paths matching any whitelisted prefix cause the walker to traferse further.

	for _, x := range wht {
		if strings.HasPrefix(pat, x) {
			return true
		}
	}

	// None of the whitelisted paths matched, so this path is not whitelisted,
	// meaning we should skip it entirely.

	return false
}
