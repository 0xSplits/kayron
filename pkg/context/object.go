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

	ind int
	kin kind
}
