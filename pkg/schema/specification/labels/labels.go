package labels

type Labels struct {
	// Environment is the name of the .yaml file that this deployment
	// configuration was found in.
	Environment string

	// Testing is true if the parent folder that this .yaml file is a child of, is
	// called "testing".
	Testing bool
}
