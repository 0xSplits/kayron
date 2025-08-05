package resolver

import (
	"context"

	"github.com/xh3b4sd/tracer"
)

func (r *Resolver) Commit() (string, error) {
	com, _, err := r.git.Repositories.GetCommit(context.Background(), r.own, r.rep, "HEAD", nil)
	if err != nil {
		return "", tracer.Mask(err)
	}

	return com.GetSHA(), nil
}
