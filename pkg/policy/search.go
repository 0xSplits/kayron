package policy

import "github.com/0xSplits/kayron/pkg/cache"

func (p *Policy) Search() []cache.Object {
	var lis []cache.Object

	for _, x := range p.cac.Releases() {
		if drift(x) {
			lis = append(lis, x)
		}
	}

	return lis
}
