package operator

import (
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/google/go-github/v76/github"
)

func newGit(env envvar.Env) *github.Client {
	// The private key option for creating github clients would look something
	// like this.
	//
	//     import "github.com/bradleyfalzon/ghinstallation/v2"
	//     import "github.com/google/go-github/v76/github"
	//
	//     key, err := base64.StdEncoding.DecodeString(env.GithubPrivateKey)
	//     if err != nil {
	//       tracer.Panic(tracer.Mask(err))
	//     }
	//
	//     itr, err := ghinstallation.New(http.DefaultTransport, env.GithubAppId, GithubInstallationId, key)
	//     if err != nil {
	//       tracer.Panic(tracer.Mask(err))
	//     }
	//
	//     return github.NewClient(&http.Client{Transport: itr})
	//

	return github.NewClient(nil).WithAuthToken(env.GithubToken)
}
