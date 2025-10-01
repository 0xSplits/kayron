package release

import (
	"fmt"

	"github.com/distribution/reference"
	"github.com/xh3b4sd/tracer"
	"golang.org/x/mod/semver"
)

// Release must be a semver version string representing the respective Github
// release tag. The tag format is required to contain a leading "v" prefix and
// optional "-" separator for providing additional metadata.
//
//	v0.1.0            the very first development release for new projects
//	v1.8.2            the fully qualified first major release for stable APIs
//	v1.8.3-ffce1e2    the metadata version for third party projects like Alloy
type String string

func (s String) Empty() bool {
	return s == ""
}

func (s String) String() string {
	return string(s)
}

func (s String) Verify() error {
	var err error

	// Ensure the shape of our release tags complies with our desired format.
	//
	//     [v.MAJOR.MINOR.PATCH(-SUFFIX)]
	//

	if !semver.IsValid(string(s)) {
		return tracer.Mask(invalidReleaseFormatError, tracer.Context{Key: "tag", Value: s})
	}

	// Explicitely validate against meta data suffixes. The + character is not
	// allowed to be part of a container image tag.
	//
	//     https://pkg.go.dev/github.com/containers/image/v5/docker/reference
	//

	if semver.Build(string(s)) != "" {
		return tracer.Mask(invalidReleaseFormatError, tracer.Context{Key: "tag", Value: s})
	}

	// Douple check that the given release tag can be used as container image tag.

	var ref reference.Reference
	{
		ref, err = reference.ParseAnyReference(fmt.Sprintf("registry/repository:%s", s))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if _, typ := ref.(reference.Tagged); !typ {
		return tracer.Mask(invalidReleaseFormatError, tracer.Context{Key: "tag", Value: s})
	}

	return nil
}
