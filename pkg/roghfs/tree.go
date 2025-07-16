package roghfs

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/tracer"
)

func (r *Roghfs) tree() error {
	var err error

	// Snynchronize calls to tree() exclusively so that we control the
	// initialization path.
	{
		r.mut.Lock()
	}

	// If we already initialized, unlock and return early. The only performance
	// cost for us is to lock and unlock on every file read, which, given the
	// problem domain that we are solving for, should not be an issue.
	if r.ini {
		r.mut.Unlock()
		return nil
	}

	// Defer the unlock if we are initializing. This code path is only executed
	// once, so that the extra cost of defer() is not accumulating during normal
	// use after initialization.
	{
		defer r.mut.Unlock()
	}

	// Get the tree structure of the configured remote Github repository
	// recursively in a single network call. Note that this limit for the tree
	// array is 100,000 entries with a maximum size of 7 MB when using the
	// recursive parameter.
	//
	//     https://docs.github.com/en/rest/git/trees?apiVersion=2022-11-28#get-a-tree
	//

	var tre *github.Tree
	{
		tre, _, err = r.git.Git.GetTree(context.Background(), r.org, r.rep, r.ref, true)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for _, x := range tre.Entries {
		var pat string
		{
			pat = filepath.Clean(x.GetPath())
		}

		// Create new directory in the underlying base file system if Github's entry
		// type is "tree".

		if x.GetType() == "tree" {
			err = r.bas.MkdirAll(pat, os.ModePerm)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		// Create new file in the underlying base file system if Github's entry type
		// is "blob". Providing nil bytes to write() will create an empty file
		// without any content.

		if x.GetType() == "blob" {
			err = r.write(pat, nil)
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	// Signal the completeness of initialization internally, so that consecutive
	// calls of tree() do not make network calls anymore.
	{
		r.ini = true
	}

	return nil
}
