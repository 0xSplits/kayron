package branch

// String triggers branch-based deployment for test environments. For instance,
// it is possible to instruct Kayron to deploy Specta according to all changes
// of the specified feature branch in the Specta repository, and deploy those
// changes to the test environment matching the current branch. Note that this
// overwrite is only considered valid for existing test environments.
type String string

func (s String) Empty() bool {
	return s == ""
}

func (s String) String() string {
	return string(s)
}

func (s String) Verify() error {
	// TODO
	return nil
}
