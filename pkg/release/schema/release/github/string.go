package github

import (
	"strings"

	"github.com/xh3b4sd/tracer"
)

type String string

func (s String) Empty() bool {
	return s == ""
}

func (s String) String() string {
	return string(s)
}

func (s String) Verify() error {
	if s == "" {
		return tracer.Mask(invalidGithubRepositoryError,
			tracer.Context{Key: "reason", Value: "repository must not be empty"},
		)
	}

	if len(s) > 100 {
		return tracer.Mask(invalidGithubRepositoryError,
			tracer.Context{Key: "reason", Value: "repository must not have more than 100 characters"},
			tracer.Context{Key: "input", Value: s},
		)
	}

	if strings.HasPrefix(string(s), ".") {
		return tracer.Mask(invalidGithubRepositoryError,
			tracer.Context{Key: "reason", Value: "repository must not start with punctuation"},
			tracer.Context{Key: "input", Value: s},
		)
	}

	if strings.HasSuffix(string(s), ".git") {
		return tracer.Mask(invalidGithubRepositoryError,
			tracer.Context{Key: "reason", Value: "repository must not end with .git"},
			tracer.Context{Key: "input", Value: s},
		)
	}

	for _, c := range s {
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' && c != '_' && c != '.' {
			return tracer.Mask(invalidGithubRepositoryError,
				tracer.Context{Key: "reason", Value: "repository must only contain A-Z a-z 0-9 - _ ."},
				tracer.Context{Key: "input", Value: s},
			)
		}
	}

	return nil
}
