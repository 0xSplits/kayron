package webhook

type Webhook []string

func (w Webhook) Empty() bool {
	return len(w) == 0
}

func (w Webhook) Verify() error {
	return nil
}
