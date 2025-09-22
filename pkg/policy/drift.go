package policy

import "github.com/0xSplits/kayron/pkg/cache"

func (p *Policy) Drift() []cache.Object {
	var lis []cache.Object

	for _, x := range p.cac.Releases() {
		if drift(x) {
			lis = append(lis, x)
		}
	}

	return lis
}

// drift tries to detect a single valid state drift, in order to allow
// allow the operator chain to execute. Our
// current policy requires the following conditions to be true for a valid
// state drift.
//
//  1. the desired deployment must not be suspended
//
//  2. the current and desired state must not be equal
//
//  3. the desired state must not be empty
//
//  4. the container image for the desired state must be pushed
func drift(obj cache.Object) bool {
	return !bool(obj.Release.Deploy.Suspend) &&
		obj.Artifact.Drift() &&
		!obj.Artifact.Empty() &&
		obj.Artifact.Valid()
}
