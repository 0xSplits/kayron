package branch

type Branch string

func (b Branch) Empty() bool {
	return b == ""
}

func (b Branch) Verify() error {
	return nil
}
