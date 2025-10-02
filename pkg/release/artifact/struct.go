// Package artifact defines the structural details of any given release
// artifact, so that the respective current and desired state can be mapped to
// their source code.
package artifact

import (
	"github.com/0xSplits/kayron/pkg/release/artifact/condition"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/artifact/scheduler"
)

type Struct struct {
	// Condition is the namespace for the state requirements of any given
	// infrastructure or service release.
	Condition condition.Struct

	// Reference is the namespace for the desired state of any given
	// infrastructure or service release.
	Reference reference.Struct

	// Scheduler is the namespace for the current state of any given
	// infrastructure or service release.
	Scheduler scheduler.Struct
}

// Drift returns whether the current state is different from the desired state,
// or whether this artifcat should trigger a deployment regardless.
func (s Struct) Drift() bool {
	return s.Scheduler.Current != s.Reference.Desired || s.Condition.Trigger
}

// Empty returns whether the desired reference of this artifact is empty.
func (s Struct) Empty() bool {
	return s.Reference.Desired == ""
}

// Merge applies the given forward only patch, which means that only non-zero
// values can overwrite zero values. This particular patch strategy is important
// because Kayron manages release artifacts concurrently, so we must only ever
// update the leafs of an artifact where it actually changed. Otherwise zero
// values would overwrite artifact leafs that are patched by another goroutine
// in parallel.
func (s Struct) Merge(p Struct) Struct {
	if s.Condition.Success == false && p.Condition.Success != false { // nolint:gosimple
		s.Condition.Success = p.Condition.Success
	}

	if s.Condition.Trigger == false && p.Condition.Trigger != false { // nolint:gosimple
		s.Condition.Trigger = p.Condition.Trigger
	}

	if s.Reference.Desired == "" && p.Reference.Desired != "" {
		s.Reference.Desired = p.Reference.Desired
	}

	if s.Scheduler.Current == "" && p.Scheduler.Current != "" {
		s.Scheduler.Current = p.Scheduler.Current
	}

	return s
}

func (s Struct) Valid() bool {
	return s.Condition.Success
}
