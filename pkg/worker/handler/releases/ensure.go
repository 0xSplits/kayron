package releases

import (
	"github.com/0xSplits/kayron/pkg/release/loader"
	"github.com/0xSplits/kayron/pkg/release/schema"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/0xSplits/kayron/pkg/worker/handler/releases/resolver"
	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Ensure() error {
	var err error

	// Figure out which Git ref to look at when fetching release information. See
	// the resolver documentation for the rules applied per environment. The
	// current implementation supports multiple test environments.

	var ref string
	{
		ref, err = resolver.Search(h.res, h.env.Environment)
		if err != nil {
			return tracer.Mask(err)
		}
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
			Git: h.git,
			Own: h.own,
			Rep: h.rep,
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

	// Cache all configured releases so that the cache key is the service index of
	// our schema list. That way other worker handlers can iterate over all cached
	// release settings like shown below. The first return value is the service
	// definition containing the respective release information. And the second
	// return value indicates whether the cache key exists, which should always be
	// true if any service release was cached in the first place.
	//
	//     for i := range cache.Length() {
	//       ser, exi := cache.Search(i)
	//     }
	//

	for i, x := range sch.Service {
		{
			h.rel.Create(i, x)
		}

		h.log.Log(
			"level", "debug",
			"message", "cached service release",
			"github", x.Github.String(),
			"deploy", x.Deploy.String(),
		)
	}

	return nil
}
