// Package roghfs implements a read-only afero.Fs for remote Github
// repositories. If you do not trust the read-only guarantees of this
// implementation, then you can wrap it in afero's own read-only interface via
// afero.NewReadOnlyFs(). Roghfs fetches the remote source files from the
// configured remote Github repository on first file system read, and delegates
// all further I/O operations to the injected base file system, e.g.
// afero.NewMemMapFs().
package roghfs

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/google/go-github/v73/github"
	"github.com/spf13/afero"
	"github.com/xh3b4sd/choreo/success"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Bas is the interface for the base file system that all repository source
	// files are written to.
	Bas afero.Fs

	// Git is the authenticated Github client used to access the configured
	// repository source files.
	Git *github.Client

	// Own is the name of the Github organization that owns the repository to read
	// from.
	Own string

	// Rep is the name of the Github repository to read from.
	Rep string

	// Ref is the Git specific branch, tag, or commit. The reserved value "HEAD"
	// must not be provided.
	Ref string
}

type Roghfs struct {
	bas afero.Fs
	git *github.Client
	own string
	rep string
	ref string

	// cac is the internal donwload cache telling us which source files we have
	// alreadz fetched. This cache is necessary because we are initializing the
	// configured root directory inside the injected base file system with empty
	// files using a single network call to Github's Tree API. This is most
	// efficient to minimize rate limit errors, but implies that we cannot tell
	// actually empty files from those that we actually have to fetch content for.
	// So this cache tells us which files we already downloaded.
	cac cache.Interface[string, struct{}]

	// mut is a concurrency helper used to synchronize the initialization of the
	// entire repository file structure inside the injected base file system, so
	// that we can ensure to only call the Github API exactly one time for that
	// particular setup task.
	mut *success.Mutex
}

func New(c Config) *Roghfs {
	if c.Bas == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Bas must not be empty", c)))
	}
	if c.Git == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Git must not be empty", c)))
	}
	if c.Own == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Own must not be empty", c)))
	}
	if c.Rep == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Rep must not be empty", c)))
	}
	if c.Ref == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ref must not be empty", c)))
	}
	if c.Ref == "HEAD" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ref must not be HEAD", c)))
	}

	return &Roghfs{
		bas: c.Bas,
		git: c.Git,
		own: c.Own,
		rep: c.Rep,
		ref: c.Ref,

		cac: cache.New[string, struct{}](),
		mut: success.New(success.Config{}),
	}
}
