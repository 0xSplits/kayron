//go:build integration

package operator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/0xSplits/kayron/pkg/webhook"
	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/workit/registry"
	"github.com/0xSplits/workit/worker/sequence"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/logger"
	"go.opentelemetry.io/otel/metric"
)

// Test_Operator_Integration runs the entire operator chain against all network
// dependencies, as if the operator was deploying for real. The required Github
// auth token to run this integration test needs at least public repo
// permissions. Also, a set of standard AWS credentials is required with
// read-only access as described in the README.md.
//
//	KAYRON_GITHUB_TOKEN=todo go test -tags=integration ./pkg/operator -v -race -run Test_Operator_Integration
func Test_Operator_Integration(t *testing.T) {
	var env envvar.Env

	{
		env = envvar.Env{
			CloudformationStack: "server-test",
			Environment:         "testing",
			GithubToken:         envvar.MustGithub(),
			LogLevel:            "debug",
			ReleaseSource:       "https://github.com/0xSplits/releases",
			S3Bucket:            "splits-cf-templates",
		}
	}

	var log logger.Interface
	{
		log = logger.New(logger.Config{
			Filter: logger.NewLevelFilter(env.LogLevel),
			Format: logger.JSONIndenter,
		})
	}

	var cfg aws.Config
	{
		cfg = envvar.MustAws()
	}

	var cac *cache.Cache
	{
		cac = cache.New(cache.Config{
			Log: log,
		})
	}

	var met metric.Meter
	{
		met = recorder.NewMeter(recorder.MeterConfig{
			Env: env.Environment,
			Sco: "kayron",
			Ver: runtime.Tag(),
		})
	}

	var pol *policy.Policy
	{
		pol = policy.New(policy.Config{
			Aws: cfg,
			Cac: cac,
			Env: env,
			Log: log,
		})
	}

	var whk *webhook.Webhook
	{
		whk = webhook.New(webhook.Config{
			Log: log,
		})
	}

	{
		emiEve(whk) // TODO
	}

	var ope *Operator
	{
		ope = New(Config{
			Aws: cfg,
			Cac: cac,
			Dry: true, // dry run, read only
			Env: env,
			Log: log,
			Met: met,
			Pol: pol,
			Whk: whk,
		})
	}

	var reg *registry.Registry
	{
		reg = registry.New(registry.Config{
			Env: env.Environment,
			Log: log,
			Met: met,
		})
	}

	var wor *sequence.Worker
	{
		wor = sequence.New(sequence.Config{
			Han: ope.Chain(),
			Log: log,
			Reg: reg,
		})
	}

	{
		err := wor.Ensure()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}
}

func emiEve(whk *webhook.Webhook) {
	var eve github.PushEvent

	err := json.Unmarshal([]byte(jsn), &eve)
	if err != nil {
		panic(err)
	}

	whk.Update(context.Background(), "", "", &eve)
}

const jsn = `{
  "ref": "refs/heads/fancy-feature-branch",
  "before": "0fd3a841c59782f8d63fa5d1b62611a2ece0fd48",
  "after": "0x123456789",
  "repository": {
    "id": 878586400,
    "node_id": "R_kgDONF4qIA",
    "name": "splits-lite",
    "full_name": "0xSplits/splits-lite",
    "private": false,
    "owner": {
      "name": "0xSplits",
      "email": null,
      "login": "0xSplits",
      "id": 91336227,
      "node_id": "MDEyOk9yZ2FuaXphdGlvbjkxMzM2MjI3",
      "avatar_url": "https://avatars.githubusercontent.com/u/91336227?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/0xSplits",
      "html_url": "https://github.com/0xSplits",
      "followers_url": "https://api.github.com/users/0xSplits/followers",
      "following_url": "https://api.github.com/users/0xSplits/following{/other_user}",
      "gists_url": "https://api.github.com/users/0xSplits/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/0xSplits/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/0xSplits/subscriptions",
      "organizations_url": "https://api.github.com/users/0xSplits/orgs",
      "repos_url": "https://api.github.com/users/0xSplits/repos",
      "events_url": "https://api.github.com/users/0xSplits/events{/privacy}",
      "received_events_url": "https://api.github.com/users/0xSplits/received_events",
      "type": "Organization",
      "user_view_type": "public",
      "site_admin": false
    },
    "html_url": "https://github.com/0xSplits/splits-lite",
    "description": "A minimal app for creating and distributing Splits",
    "fork": false,
    "url": "https://api.github.com/repos/0xSplits/splits-lite",
    "forks_url": "https://api.github.com/repos/0xSplits/splits-lite/forks",
    "keys_url": "https://api.github.com/repos/0xSplits/splits-lite/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/0xSplits/splits-lite/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/0xSplits/splits-lite/teams",
    "hooks_url": "https://api.github.com/repos/0xSplits/splits-lite/hooks",
    "issue_events_url": "https://api.github.com/repos/0xSplits/splits-lite/issues/events{/number}",
    "events_url": "https://api.github.com/repos/0xSplits/splits-lite/events",
    "assignees_url": "https://api.github.com/repos/0xSplits/splits-lite/assignees{/user}",
    "branches_url": "https://api.github.com/repos/0xSplits/splits-lite/branches{/branch}",
    "tags_url": "https://api.github.com/repos/0xSplits/splits-lite/tags",
    "blobs_url": "https://api.github.com/repos/0xSplits/splits-lite/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/0xSplits/splits-lite/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/0xSplits/splits-lite/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/0xSplits/splits-lite/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/0xSplits/splits-lite/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/0xSplits/splits-lite/languages",
    "stargazers_url": "https://api.github.com/repos/0xSplits/splits-lite/stargazers",
    "contributors_url": "https://api.github.com/repos/0xSplits/splits-lite/contributors",
    "subscribers_url": "https://api.github.com/repos/0xSplits/splits-lite/subscribers",
    "subscription_url": "https://api.github.com/repos/0xSplits/splits-lite/subscription",
    "commits_url": "https://api.github.com/repos/0xSplits/splits-lite/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/0xSplits/splits-lite/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/0xSplits/splits-lite/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/0xSplits/splits-lite/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/0xSplits/splits-lite/contents/{+path}",
    "compare_url": "https://api.github.com/repos/0xSplits/splits-lite/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/0xSplits/splits-lite/merges",
    "archive_url": "https://api.github.com/repos/0xSplits/splits-lite/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/0xSplits/splits-lite/downloads",
    "issues_url": "https://api.github.com/repos/0xSplits/splits-lite/issues{/number}",
    "pulls_url": "https://api.github.com/repos/0xSplits/splits-lite/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/0xSplits/splits-lite/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/0xSplits/splits-lite/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/0xSplits/splits-lite/labels{/name}",
    "releases_url": "https://api.github.com/repos/0xSplits/splits-lite/releases{/id}",
    "deployments_url": "https://api.github.com/repos/0xSplits/splits-lite/deployments",
    "created_at": 1729875783,
    "updated_at": "2025-09-09T10:07:29Z",
    "pushed_at": 1759235302,
    "git_url": "git://github.com/0xSplits/splits-lite.git",
    "ssh_url": "git@github.com:0xSplits/splits-lite.git",
    "clone_url": "https://github.com/0xSplits/splits-lite.git",
    "svn_url": "https://github.com/0xSplits/splits-lite",
    "homepage": "https://lite.splits.org",
    "size": 1462,
    "stargazers_count": 2,
    "watchers_count": 2,
    "language": "TypeScript",
    "has_issues": true,
    "has_projects": false,
    "has_downloads": true,
    "has_wiki": false,
    "has_pages": false,
    "has_discussions": false,
    "forks_count": 1,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 13,
    "license": {
      "key": "gpl-3.0",
      "name": "GNU General Public License v3.0",
      "spdx_id": "GPL-3.0",
      "url": "https://api.github.com/licenses/gpl-3.0",
      "node_id": "MDc6TGljZW5zZTk="
    },
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [],
    "visibility": "public",
    "forks": 1,
    "open_issues": 13,
    "watchers": 2,
    "default_branch": "main",
    "stargazers": 2,
    "master_branch": "main",
    "organization": "0xSplits",
    "custom_properties": {}
  },
  "pusher": {
    "name": "xh3b4sd",
    "email": "xh3b4sd@gmail.com"
  },
  "organization": {
    "login": "0xSplits",
    "id": 91336227,
    "node_id": "MDEyOk9yZ2FuaXphdGlvbjkxMzM2MjI3",
    "url": "https://api.github.com/orgs/0xSplits",
    "repos_url": "https://api.github.com/orgs/0xSplits/repos",
    "events_url": "https://api.github.com/orgs/0xSplits/events",
    "hooks_url": "https://api.github.com/orgs/0xSplits/hooks",
    "issues_url": "https://api.github.com/orgs/0xSplits/issues",
    "members_url": "https://api.github.com/orgs/0xSplits/members{/member}",
    "public_members_url": "https://api.github.com/orgs/0xSplits/public_members{/member}",
    "avatar_url": "https://avatars.githubusercontent.com/u/91336227?v=4",
    "description": "Financial infrastructure for onchain startups"
  },
  "sender": {
    "login": "xh3b4sd",
    "id": 552769,
    "node_id": "MDQ6VXNlcjU1Mjc2OQ==",
    "avatar_url": "https://avatars.githubusercontent.com/u/552769?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/xh3b4sd",
    "html_url": "https://github.com/xh3b4sd",
    "followers_url": "https://api.github.com/users/xh3b4sd/followers",
    "following_url": "https://api.github.com/users/xh3b4sd/following{/other_user}",
    "gists_url": "https://api.github.com/users/xh3b4sd/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/xh3b4sd/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/xh3b4sd/subscriptions",
    "organizations_url": "https://api.github.com/users/xh3b4sd/orgs",
    "repos_url": "https://api.github.com/users/xh3b4sd/repos",
    "events_url": "https://api.github.com/users/xh3b4sd/events{/privacy}",
    "received_events_url": "https://api.github.com/users/xh3b4sd/received_events",
    "type": "User",
    "user_view_type": "public",
    "site_admin": false
  },
  "installation": {
    "id": 87614107,
    "node_id": "MDIzOkludGVncmF0aW9uSW5zdGFsbGF0aW9uODc2MTQxMDc="
  },
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/0xSplits/splits-lite/compare/0fd3a841c597...547fd013ac2d",
  "commits": [
    {
      "id": "547fd013ac2d5b00514ed3af022f67e522605505",
      "tree_id": "8b12906e5f88c362ae6e9bec2c359b873c84e60b",
      "distinct": true,
      "message": "red",
      "timestamp": "2025-09-30T14:28:20+02:00",
      "url": "https://github.com/0xSplits/splits-lite/commit/547fd013ac2d5b00514ed3af022f67e522605505",
      "author": {
        "name": "xh3b4sd",
        "email": "xh3b4sd@gmail.com",
        "username": "xh3b4sd"
      },
      "committer": {
        "name": "xh3b4sd",
        "email": "xh3b4sd@gmail.com",
        "username": "xh3b4sd"
      },
      "added": [],
      "removed": [],
      "modified": [
        "src/app/layout.tsx"
      ]
    }
  ],
  "head_commit": {
    "id": "547fd013ac2d5b00514ed3af022f67e522605505",
    "tree_id": "8b12906e5f88c362ae6e9bec2c359b873c84e60b",
    "distinct": true,
    "message": "red",
    "timestamp": "2025-09-30T14:32:20+02:00",
    "url": "https://github.com/0xSplits/splits-lite/commit/547fd013ac2d5b00514ed3af022f67e522605505",
    "author": {
      "name": "xh3b4sd",
      "email": "xh3b4sd@gmail.com",
      "username": "xh3b4sd"
    },
    "committer": {
      "name": "xh3b4sd",
      "email": "xh3b4sd@gmail.com",
      "username": "xh3b4sd"
    },
    "added": [],
    "removed": [],
    "modified": [
      "src/app/layout.tsx"
    ]
  }
}`
