package webhook

type Webhook string

func (w Webhook) Empty() bool {
	return w == ""
}

func (w Webhook) Verify() error {
	// TODO
	return nil
}
