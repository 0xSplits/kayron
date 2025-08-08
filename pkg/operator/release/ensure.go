package release

import (
	"github.com/0xSplits/kayron/pkg/operator/release/resolver"
	"github.com/0xSplits/kayron/pkg/release/loader"
	"github.com/0xSplits/kayron/pkg/release/schema"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

func (r *Release) Ensure() error {
	var err error

	// The release handler is the very first building block in our operator chain.
	// We have to make sure that every reconciliation loop starts with a blank
	// slate. So before doing anything else, we have to purge all cache state by
	// calling Delete below.

	{
		r.cac.Delete()
	}

	// Figure out which Git ref to look at when fetching release information. See
	// the resolver documentation for the rules applied per environment. The
	// current implementation supports multiple test environments.

	var ref string
	{
		ref, err = resolver.Search(r.res, r.env.Environment)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		r.log.Log(
			"level", "debug",
			"message", "resolved ref for github repository",
			"environment", r.env.Environment,
			"repository", r.env.ReleaseSource,
			"ref", ref,
		)
	}

	// On every loop, create a new Read Only GitHub File System, to fetch the
	// latest version of the releases repository. It is important to create a new
	// base file system every time we want to refresh our view of the watched
	// repository, because that base file system reflects all repository state
	// that we know about.

	var gfs *roghfs.Roghfs
	{
		gfs = roghfs.New(roghfs.Config{
			Bas: afero.NewMemMapFs(),
			Git: r.git,
			Own: r.own,
			Rep: r.rep,
			Ref: ref,
		})
	}

	// Fetch the release settings from the configured Github repository by using
	// our standard schema loader. The behaviour of the loader is standardized, so
	// that loading from a local file system and loading from Github remains
	// logically the same.

	var sch schema.Schema
	{
		sch, err = loader.Loader(gfs, ".")
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		err = sch.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Initialize the cache for all configured releases regardless of their type.
	// Here we require exactly one infrastructure release to be provided.

	{
		err = r.cac.Create(sch.Release)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
