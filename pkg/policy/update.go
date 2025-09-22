package policy

import (
	"slices"
)

// Update determines whether the managed CloudFormation stack should be updated
// or not. We do not signal an update if the managed CloudFormation stack is
// already being updated, and if there is no detectable state drift.
func (p *Policy) Update() bool {
	// Fetch the deployment status of the underlying root stack so that we can
	// decide whether to proceed with the execution of writing operator functions.
	var can bool
	{
		can = p.Cancel()
	}

	// The slices.ContainsFunc version below is the equivalent of the shown for
	// loop.
	//
	//     for _, x := range p.cac.Releases() {
	//        if drift(x) {
	//          return true
	//        }
	//     }
	//
	var dft bool
	{
		dft = slices.ContainsFunc(p.cac.Releases(), drift)
	}

	return !can && dft
}
