package resolver

import (
	"context"

	"github.com/xh3b4sd/tracer"
)

func (r *Resolver) Latest() (string, error) {
	rel, res, err := r.git.Repositories.GetLatestRelease(context.Background(), r.own, r.rep)
	if isNotFound(res) {
		return "", tracer.Mask(releaseNotFoundError,
			tracer.Context{Key: "owner", Value: r.own},
			tracer.Context{Key: "repository", Value: r.rep},
		)
	} else if err != nil {
		return "", tracer.Mask(err)
	}

	return rel.GetTagName(), nil
}
