package infrastructure

// Infrastructure defines optional, environment specific overwrites.
type Infrastructure struct {
	// Overwrite supports infrastructure specific branch overrides for test
	// environments. Maintaining a feature branch with the name of the test
	// environment to overwrite, allows to deploy feature branch specific changes
	// automatically. E.g. creating the feature branch "testing" and modifying the
	// testing.yaml file in the folder environment/testing/, defines
	// infrastructure overwrites that are automatically reconciled within the
	// "testing" environment. Note that this feature is only available for test
	// environments and is constrained to specifically named feature branches
	// only.
	Overwrite bool `yaml:"overwrite,omitempty"`

	// Shorthand specifies an alternative version of the environment name for
	// internal resource identification. E.g. the shorthand CloudFormation stack
	// identifier for the "staging" environment is called "master" for legacy
	// reasons.
	Shorthand string `yaml:"shorthand,omitempty"`
}
