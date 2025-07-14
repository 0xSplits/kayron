package suspend

type Suspend bool

func (s Suspend) Empty() bool {
	return !bool(s)
}

func (s Suspend) Verify() error {
	return nil
}
