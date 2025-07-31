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

	// Get the tree structure of the configured remote Github repository
	// recursively in a single network call. Note that this limit for the tree
	// array is 100,000 entries with a maximum size of 7 MB when using the
	// recursive parameter.
	//
	//     https://docs.github.com/en/rest/git/trees?apiVersion=2022-11-28#get-a-tree
	//
	var tre *github.Tree
	{
		tre, _, err = r.git.Git.GetTree(context.Background(), r.own, r.rep, r.ref, true)
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

	return nil
}
