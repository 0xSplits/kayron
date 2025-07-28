package resolver

import (
	"context"
	"net/http"

	"github.com/xh3b4sd/tracer"
)

func (r *Resolver) Latest() (string, error) {
	rel, res, err := r.git.Repositories.GetLatestRelease(context.Background(), r.own, r.rep)
	if err != nil {
		if res != nil && res.StatusCode == http.StatusNotFound {
			return "", tracer.Mask(releaseNotFoundError,
				tracer.Context{Key: "owner", Value: r.own},
				tracer.Context{Key: "repo", Value: r.rep},
			)
		}
		return "", tracer.Mask(err)
	}

	return rel.GetTagName(), nil
}
