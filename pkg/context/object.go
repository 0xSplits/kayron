package context

import (
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
)

type kind string

const (
	Infrastructure kind = "infrastructure"
	Service        kind = "service"
)

type Object struct {
	Artifact artifact.Struct
	Release  release.Struct
	Kind     kind
}

func (o Object) Drift() bool {
	return o.Artifact.Scheduler.Current != o.Artifact.Reference.Desired
}

// Empty returns whether the desired reference of this artifact is empty.
func (o Object) Empty() bool {
	return o.Artifact.Reference.Desired == ""
}

func (o Object) Valid() bool {
	return o.Artifact.Condition.Success
}
