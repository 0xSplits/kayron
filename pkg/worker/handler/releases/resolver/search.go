package resolver

import (
	"github.com/xh3b4sd/tracer"
)

// Search determines the Git ref to use for the given environment. An empty ref
// indicates to use the default branch, which is in line with the documented
// behaviour of the official Github API.
//
//	production    use the "production" branch if it exists, otherwise fall back
//	              to the latest release tag
//
//	staging       always use the default branch of the underlying Github
//	              repository by returning an empty ref
//
//	testing       use the test branch matching the test environment name if it
//	              exists, otherwise fall back to the default branch
func Search(res Interface, env string) (string, error) {
	// Always use the default branch for the "staging" environment.

	if env == "staging" {
		return "", nil
	}

	// Use the existing branches for either the "production" or test environments,
	// if those branches exist.

	{
		exi, err := res.Exists(env)
		if err != nil {
			return "", tracer.Mask(err)
		}
		if exi {
			return env, nil
		}
	}

	// If no "production" branch exists, then fall back to the Git tag of the
	// latest Github release.

	if env == "production" {
		tag, err := res.Latest()
		if err != nil {
			return "", tracer.Mask(err)
		}

		return tag, nil
	}

	// Use the default branch if no branch exist for the given test environment.

	return "", nil
}
