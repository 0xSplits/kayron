package branch

// Branch triggers branch-based deployment for test environments. For instance,
// it is possible to instruct Kayron to deploy Specta according to all changes
// of the specified feature branch in the Specta repository, and deploy those
// changes to the test environment matching the current branch. Note that this
// overwrite is only considered valid for existing test environments.
type Branch string

func (b Branch) Empty() bool {
	return b == ""
}

func (b Branch) Verify() error {
	// TODO
	return nil
}
