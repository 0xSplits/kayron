package labels

// Labels are runtime specific internals annotated inside the schema loader.
type Labels struct {
	// Environment is the name of the .yaml file that this deployment
	// configuration was found in.
	Environment string

	// Source is the configuration's source file path as loaded from the
	// underlying file system. This label may help to make error messages more
	// useful.
	Source string

	// Testing is true if the parent folder that this .yaml file is a child of, is
	// called "testing".
	Testing bool
}
