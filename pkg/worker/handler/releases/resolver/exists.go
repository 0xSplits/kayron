package resolver

import (
	"context"
	"net/http"

	"github.com/xh3b4sd/tracer"
)

func (r *Resolver) Exists(bra string) (bool, error) {
	_, res, err := r.git.Repositories.GetBranch(context.Background(), r.own, r.rep, bra, 3)
	if err != nil {
		if res != nil && res.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, tracer.Mask(err)
	}

	return true, nil
}
