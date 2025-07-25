package resolver

import (
	"context"

	"github.com/xh3b4sd/tracer"
)

func (r *Resolver) Latest() (string, error) {
	rel, _, err := r.git.Repositories.GetLatestRelease(context.Background(), r.own, r.rep)
	if err != nil {
		return "", tracer.Mask(err)
	}

	if rel == nil || rel.GetTagName() == "" {
		return "", tracer.Mask(releaseNotFoundError)
	}

	return rel.GetTagName(), nil
}
