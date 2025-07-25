package roghfs

import (
	"net/url"
	"strings"

	"github.com/xh3b4sd/tracer"
)

// Parse tries to return the repository owner and the repository name for a
// standard GitHub repository URL.
//
//	https://github.com/owner/repo
func Parse(str string) (string, string, error) {
	var err error

	var trm string
	{
		trm = strings.TrimSpace(str)
		if str == "" {
			return "", "", tracer.Mask(invalidRepositoryUrlError, tracer.Context{Key: "reason", Value: "url must not be empty"})
		}
	}

	var prs *url.URL
	{
		prs, err = url.Parse(trm)
		if err != nil {
			return "", "", tracer.Mask(err)
		}
	}

	var pth string
	var spl []string
	{
		pth = strings.Trim(prs.Path, "/")
		spl = strings.Split(pth, "/")
	}

	{
		if prs.Scheme != "https" {
			return "", "", tracer.Mask(invalidRepositoryUrlError,
				tracer.Context{Key: "reason", Value: "scheme must be https://"},
				tracer.Context{Key: "received", Value: prs.Scheme},
			)
		}
		if prs.Host != "github.com" {
			return "", "", tracer.Mask(invalidRepositoryUrlError,
				tracer.Context{Key: "reason", Value: "host must be github.com"},
				tracer.Context{Key: "received", Value: prs.Host},
			)
		}
		if prs.RawQuery != "" {
			return "", "", tracer.Mask(invalidRepositoryUrlError,
				tracer.Context{Key: "reason", Value: "query must be empty"},
				tracer.Context{Key: "received", Value: prs.RawQuery},
			)
		}
		if prs.Fragment != "" {
			return "", "", tracer.Mask(invalidRepositoryUrlError,
				tracer.Context{Key: "reason", Value: "fragment must be empty"},
				tracer.Context{Key: "received", Value: prs.Fragment},
			)
		}
		if len(spl) != 2 || spl[0] == "" || spl[1] == "" {
			return "", "", tracer.Mask(invalidRepositoryUrlError,
				tracer.Context{Key: "reason", Value: "format must be owner/repo"},
				tracer.Context{Key: "received", Value: spl},
			)
		}
		if strings.HasSuffix(spl[1], ".git") {
			return "", "", tracer.Mask(invalidRepositoryUrlError,
				tracer.Context{Key: "reason", Value: "repo must not contain .git suffix"},
				tracer.Context{Key: "received", Value: spl[1]},
			)
		}
	}

	return spl[0], spl[1], nil
}
