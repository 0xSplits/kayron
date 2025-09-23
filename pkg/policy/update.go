package policy

// Update determines whether the managed CloudFormation stack should be updated
// or not. We do not signal an update if the managed CloudFormation stack is
// already being updated, and if there is no detectable state drift.
func (p *Policy) Update() bool {
	// Fetch the deployment status of the underlying root stack so that we can
	// decide whether to proceed with the execution of writing operator functions.

	var can bool
	{
		can = p.Cancel()
	}

	// Figure out whether we have any state drift at all that has to be
	// reconciled.

	var dft bool
	{
		dft = p.drift()
	}

	return !can && dft
}

func (p *Policy) drift() bool {
	for _, x := range p.cac.Releases() {
		if x.Drift() {
			return true
		}
	}

	return false
}
