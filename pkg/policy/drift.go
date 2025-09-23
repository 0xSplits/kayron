package policy

import "github.com/0xSplits/kayron/pkg/cache"

// Drift returns all cached artifact releases that have valid state drift. In
// other words, the cache objects returned here indicate that their respective
// releases should be updated.
func (p *Policy) Drift() []cache.Object {
	var lis []cache.Object

	for _, x := range p.cac.Releases() {
		if x.Drift() {
			lis = append(lis, x)
		}
	}

	return lis
}
