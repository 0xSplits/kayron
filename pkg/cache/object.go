package cache

import (
	"fmt"
	"strings"

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

func (o Object) Name() string {
	if o.kin == Infrastructure {
		return o.Release.Github.String()
	}

	if o.kin == Service {
		return o.Release.Docker.String()
	}

	return ""
}

// Parameter returns the CloudFormation stack parameter key for this release
// artifact. The parameter keys generated here have to be supported in the
// CloudFormation template being deployed, e.g. InfrastructureVersion,
// SpectaVersion.
func (o Object) Parameter() string {
	return fmt.Sprintf("%sVersion", strings.Map(mapFnc, cases.Title(language.English).String(o.Name())))
}

// Version returns the desired state of this artifact's release version if the
// following two conditions are true. Failed artifact conditions and deployment
// suspensions will then yield the artifact version of the scheduler's current
// state.
//
//  1. the artifact condition must be true
//
//  2. the deployment strategy must not be suspended
//
// E.g. we may observe the new version of a service v0.5.0, for which there does
// no docker image exist yet. In this case we defer to the old version v0.4.0,
// which is currently deployed. Once the docker image for v0.5.0 has been
// confirmed inside of the underlying container registry, we will yield the new
// desired version.
func (o Object) Version() string {
	if bool(o.Release.Deploy.Suspend) || !o.Artifact.Valid() {
		return o.Artifact.Scheduler.Current
	}

	return o.Artifact.Reference.Desired
}

// mapFnc signals the removal of dashes and spaces only.
func mapFnc(r rune) rune {
	if r == '-' || r == ' ' {
		return -1
	}

	return r
}
