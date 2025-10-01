package reference

import (
	"net/http"

	"github.com/google/go-github/v75/github"
)

func isNotFound(res *github.Response) bool {
	if res == nil {
		return false
	}

	return res.StatusCode == http.StatusNotFound
}
