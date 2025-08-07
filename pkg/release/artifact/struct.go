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
