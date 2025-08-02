package release

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/release/loader"
	"github.com/0xSplits/kayron/pkg/release/schema"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/release/resolver"
	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

func (r *Release) Ensure() error {
	var err error

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

	// Create the release cache for all configured services so that the cache key
	// aligns with its respective service index across all caches. That way, the
	// following business logic can iterate over all cached artifact and release
	// settings like shown below. The first return value is the respective cache
	// value, and the second return value indicates whether the cache key exists,
	// which should always be true on idiomatic iteration.
	//
	//     for i := range r.ser.Length() {
	//       ser, _ := r.ser.Search(i)
	//     }
	//

	for i, x := range sch.Service {
		{
			r.ser.Create(i, x)
		}

		r.log.Log(
			"level", "debug",
			"message", "cached service release",
			"docker", x.Docker.String(),
			"github", x.Github.String(),
			"deploy", x.Deploy.String(),
			"artifact", fmt.Sprintf("[%d]", i),
		)
	}

	return nil
}
