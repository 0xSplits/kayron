package webhook

// Webhook contains a list of alternative deployment mechanisms. Each webhook
// provided here is invoked to deploy e.g. our frontends in Vercel. The format
// of those webhook definitions requires the usage of a prefix for the HTTP
// method that this webhook should be called with. It is further required to
// provide a HTTPs URL. Failed webhook calls may be retried and eventually be
// reported as terminal failure.
//
//	POST:https://{{DNS}}/{{PATH}}
type Webhooks []Webhook

func (w Webhooks) Empty() bool {
	return len(w) == 0
}

func (w Webhooks) Verify() error {
	// TODO
	return nil
}
