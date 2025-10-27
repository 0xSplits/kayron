package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v76/github"
	"github.com/xh3b4sd/tracer"
)

// Update implements PushEventHandleFunc. We bind this method to the webhook
// endpoint that Github is POSTing push events to.
//
//	https://docs.github.com/en/webhooks/webhook-events-and-payloads#push
//	https://pkg.go.dev/github.com/cbrgm/githubevents/v2@v2.5.0/githubevents#PushEventHandleFunc
func (w *Webhook) Update(ctx context.Context, did string, nam string, eve *github.PushEvent) error {
	// First, some safe guards. We are also not interested in branch and tag
	// deletions.

	if eve == nil || eve.GetDeleted() {
		return nil
	}

	// Ignore tags and other ref structures. We are only interested in commit
	// pushes into branches.

	var ref string
	{
		ref = eve.GetRef()
	}

	if !strings.HasPrefix(ref, branch) {
		return nil
	}

	// Parse the branch name of the push event, so we can create a cache key.

	var bra string
	{
		bra = strings.TrimPrefix(ref, branch)
	}

	// Create the cache key. We use the login of the owner, because not every
	// commit has an organization. We just call it org internally.

	var key Key
	{
		key = Key{
			Org: eve.GetRepo().GetOwner().GetLogin(),
			Rep: eve.GetRepo().GetName(),
			Bra: bra,
		}
	}

	// Create the commit object we would like to cache. We have to use the "after"
	// field of the push event payload, because this is the most reliable commit
	// hash of the head commit. The webhook API does not populate the "sha" field
	// in the head commit object. ¯\_(ツ)_/¯

	var com Commit
	{
		com = Commit{
			Hash: eve.GetAfter(),
			Time: eve.GetHeadCommit().GetTimestamp().Time,
		}
	}

	// Verify that our commit object has a non empty hash and timestamp.

	{
		err := w.verify(com, eve)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Finally cache the received commit hash, but only if the version we just
	// created is in fact the latest commit based on our internal cache state.

	var lat Commit
	{
		lat = w.Latest(key, com)
	}

	// If our new commit is the latest, then we want to store it. Also make sure
	// to synchronize cache access using a mutex.

	if com.Equals(lat) {
		w.log.Log(
			"level", "info",
			"message", "caching push event",
			"hash", com.Hash,
			"repository", fmt.Sprintf("https://github.com/%s/%s", key.Org, key.Rep),
			"ref", key.Bra,
			"timestamp", com.Time.String(),
		)

		{
			w.mut.Lock()
			w.cac[key] = com
			w.mut.Unlock()
		}
	}

	return nil
}

func (w *Webhook) verify(com Commit, eve *github.PushEvent) error {
	err := com.Verify()
	if err != nil {
		byt, jrr := json.Marshal(eve)
		if jrr != nil {
			return tracer.Mask(jrr)
		}

		w.log.Log(
			"level", "warning",
			"message", "skipping push event cache",
			"reason", err.Error(),
			"event", string(byt),
		)

		return tracer.Mask(err)
	}

	return nil
}
