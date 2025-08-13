package cache

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// Parameter returns the CloudFormation stack parameter key for this release
// artifact. The keys generated here have to be supported in the CloudFormation
// template being deployed.
func (o Object) Parameter() string {
	cas := cases.Title(language.English)

	if o.kin == Infrastructure {
		return fmt.Sprintf("%sVersion", cas.String(o.Release.Github.String())) // e.g. InfrastructureVersion
	}

	if o.kin == Service {
		return fmt.Sprintf("%sVersion", cas.String(o.Release.Docker.String())) // e.g. SpectaVersion
	}

	return ""
}
