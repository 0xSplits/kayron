package cache

import (
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
)

// kind is a private type for private object properties in order to guarantee
// differentiated access patterns for infrastructure and service releases.
type kind string

const (
	Infrastructure kind = "infrastructure"
	Service        kind = "service"
)

// Object combines the associated artifact and release information in one
// addressable cache object. The artifact information are written and the
// release information are read throughout the reconciliation loops.
type Object struct {
	Artifact artifact.Struct
	Release  release.Struct

	ind int
	kin kind
}
