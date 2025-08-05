package resolver

import (
	"context"

	"github.com/xh3b4sd/tracer"
)

func (r *Resolver) Commit(ref string) (string, error) {
	com, res, err := r.git.Repositories.GetCommit(context.Background(), r.own, r.rep, ref, nil)
	if isNotFound(res) {
		return "", nil
	} else if err != nil {
		return "", tracer.Mask(err)
	}

	return com.GetSHA(), nil
}
