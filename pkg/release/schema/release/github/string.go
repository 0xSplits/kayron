package github

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
