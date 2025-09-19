package preview

import (
	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

func (p *Preview) Ensure() error {
	// Get the list of cached releases so that we can lookup their respective
	// artifact references for any potential preview deployment settings.

	var rel []cache.Object
	{
		rel = p.cac.Releases()
	}

	fnc := func(_ int, o cache.Object) error {
		var err error

		// If this release has preview deployments disabled, then ignore this cache
		// object and move on to the next one.

		if !bool(o.Release.Deploy.Preview) {
			return nil
		}

		// If this release has preview deployments enabled, then compute the preview
		// releases, so that we can inject them into the internal cache below.

		var exp release.Slice
		{
			exp, err = p.pre.Expand(o.Release)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		// Extend the cache for all expanded preview deployments.

		{
			err := p.cac.Create(exp)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		return nil
	}

	{
		err := parallel.Slice(rel, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
