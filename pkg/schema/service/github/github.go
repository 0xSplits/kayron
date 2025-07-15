package github

type Github string

func (g Github) Empty() bool {
	return g == ""
}

func (g Github) Verify() error {
	// TODO
	return nil
}
