package resolver

import (
	"github.com/xh3b4sd/tracer"
)

// Search determines the Git ref to use for the given environment. Search will
// never return an empty ref, but instead resolve the SHA of the Git commit at
// HEAD of the default branch.
//
//	production    resolve the latest release tag
//
//	staging       resolve the latest commit sha of the default branch
//
//	testing       resolve the latest commit sha of the branch matching the test
//	              environment, if such a branch exists, otherwise fall back to
//	              the commit sha of the default branch
func Search(res Interface, env string) (string, error) {
	// Always use the commit sha of the default branch for the "staging"
	// environment.

	if env == "staging" {
		sha, err := res.Commit("HEAD")
		if err != nil {
			return "", tracer.Mask(err)
		}

		return sha, nil
	}

	// Use the Git tag of the latest Github release for the "production"
	// environment.

	if env == "production" {
		tag, err := res.Latest()
		if err != nil {
			return "", tracer.Mask(err)
		}

		return tag, nil
	}

	// Use the commmit sha of any existing branch for any given test environment,
	// if such a branch exists.

	{
		sha, err := res.Commit(env)
		if err != nil {
			return "", tracer.Mask(err)
		}

		if sha != "" {
			return sha, nil
		}
	}

	// Use the commit sha of the default branch if no branch exist for the given
	// test environment.

	{
		sha, err := res.Commit("HEAD")
		if err != nil {
			return "", tracer.Mask(err)
		}

		return sha, nil
	}
}
