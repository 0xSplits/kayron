package docker

type String string

func (s String) Empty() bool {
	return s == ""
}

func (s String) Verify() error {
	// TODO
	return nil
}
