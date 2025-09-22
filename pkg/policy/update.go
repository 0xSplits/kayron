package policy

import "slices"

func (p *Policy) Update() bool {
	// The slices.ContainsFunc version below is the equivalent of the shown for
	// loop.
	//
	//     for _, x := range p.cac.Releases() {
	//        if drift(x) {
	//          return true
	//        }
	//     }
	//
	return slices.ContainsFunc(p.cac.Releases(), drift)
}
