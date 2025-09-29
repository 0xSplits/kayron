package webhook

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/google/go-github/v73/github"
)

func (w *Webhook) Update(ctx context.Context, did string, nam string, eve *github.PushEvent) error {
	if eve == nil || eve.GetDeleted() {
		return nil
	}

	var ref string
	{
		ref = eve.GetRef()
	}

	// ignore tags etc
	if !strings.HasPrefix(ref, branch) {
		return nil
	}

	var bra string
	{
		bra = strings.TrimPrefix(ref, branch)
	}

	var hea *github.HeadCommit
	{
		hea = eve.GetHeadCommit()
	}

	// "owner/repo/branch"
	var key string
	{
		key = path.Join(eve.GetRepo().GetFullName(), bra)
	}

	var com Commit
	{
		com = Commit{
			Hash: hea.GetSHA(),
			Time: hea.GetTimestamp().Time,
		}
	}

	fmt.Printf("\n")
	fmt.Printf("Caching Webhook Event\n")
	fmt.Printf("    %#v\n", key)
	fmt.Printf("    %#v\n", com.Hash)
	fmt.Printf("    %#v\n", com.Time.String())
	fmt.Printf("\n")

	// TODO guard against empty hash/time
	// TODO only cache if head commit is newer

	{
		w.mut.Lock()
		w.cac[key] = com
		w.mut.Unlock()
	}

	return nil
}
