package docker

type Docker string

func (d Docker) Empty() bool {
	return d == ""
}

func (d Docker) Verify() error {
	// TODO
	return nil
}
