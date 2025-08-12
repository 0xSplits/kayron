package condition

type Struct struct {
	// Success indicates whether the desired state requirements for any given
	// scheduler are fulfilled. E.g. the container image tag of any given service
	// release must be pushed to the configured container registry before its
	// desired service release can be rolled out.
	Success bool

	// Trigger indicates whether the desired state should be applied regardless of
	// any detectable drift. E.g. the changes to some modified CloudFormation tags
	// may have to be applied, which may not be detectable for the operator, and
	// hence would not be automatically reconciled otherwise.
	Trigger bool
}
