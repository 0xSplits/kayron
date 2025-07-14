package release

type Release string

func (r Release) Empty() bool {
	return r == ""
}

func (r Release) Verify() error {
	return nil
}

// TODO
// func (d *Deploy) Verify() error {
// 	// A *very* small, fast check; rely on a library like
// 	// github.com/Masterminds/semver/v3 for full parsing if needed.
// 	if len(d.Release) < 2 || d.Release[0] != 'v' {
// 		return fmt.Errorf("TODO error")
// 	}
// 	_, err := semver.NewVersion(d.Release[1:])
// 	return err
// }
