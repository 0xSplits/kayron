package envvar

import (
	"os"
)

func MustGithub() string {
	var tok string
	{
		tok = os.Getenv("KAYRON_GITHUB_TOKEN")
		if tok == "" {
			panic("env var KAYRON_GITHUB_TOKEN must not be empty")
		}
	}

	return tok
}
