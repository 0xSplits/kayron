package roghfs

import (
	"context"
	"io"

	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/tracer"
)

// search downloads the content of a source file from the configured remote
// Github repository and returns its raw bytes.
func (r *Roghfs) search(pat string) ([]byte, error) {
	var err error

	var opt *github.RepositoryContentGetOptions
	{
		opt = &github.RepositoryContentGetOptions{
			Ref: r.ref,
		}
	}

	var rea io.ReadCloser
	{
		rea, _, err = r.git.Repositories.DownloadContents(context.Background(), r.own, r.rep, pat, opt)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	{
		defer rea.Close() // nolint:errcheck
	}

	var byt []byte
	{
		byt, err = io.ReadAll(rea)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return byt, nil
}
