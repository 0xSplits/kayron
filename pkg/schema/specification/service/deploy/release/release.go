package release

// Release must be a semver version string representing the respective Github
// release tag. The tag format is required to contain a leading "v" prefix and
// optional "-" separator for providing additional metadata.
//
//	v0.1.0            the very first development release for new projects
//	v1.8.2            the fully qualified first major release for stable APIs
//	v1.8.3-ffce1e2    the metadata version for third party projects like Alloy
type Release string

func (r Release) Empty() bool {
	return r == ""
}

func (r Release) Verify() error {
	// TODO
	return nil
}

// TODO
// func (d *Deploy) Verify() error {
// 	// A *very* small, fast check; rely on a library like
// 	// github.com/Masterminds/semver/v3 for full parsing if needed.
// 	if len(d.Release) < 2 || d.Release[0] != 'v' {
// 		return fmt.Errorf("custom error")
// 	}
// 	_, err := semver.NewVersion(d.Release[1:])
// 	return err
// }
