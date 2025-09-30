package webhook

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/logger"
)

func Test_Webhook_Lifecycle(t *testing.T) {
	var whk *Webhook
	{
		whk = New(Config{
			Log: logger.Fake(),
		})
	}

	var key Key
	{
		key = Key{
			Org: "a",
			Rep: "b",
			Bra: "c",
		}
	}

	// Our cache key should not exist inside an empty cache.

	var exi bool
	{
		_, exi = whk.cac[key]
		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	//
	// First Commit 0x1234
	//

	var one Commit
	{
		one = Commit{
			Hash: "0x1234",
			Time: time.Unix(3, 0),
		}
	}

	var lat Commit
	{
		lat = whk.Latest(key, one)
	}

	// The webhook cache is empty. When we ask for the latest commit using our new
	// commit object as compared version, then we should only get our compared
	// version back as the latest commit.

	if !one.Equals(lat) {
		t.Fatal("expected", one.Hash, "got", lat.Hash)
	}

	// Ensure that Webhook.Latest does not store anything in the cache.

	{
		_, exi = whk.cac[key]
		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	// Put the commit into the internal cache by way of calling our webhook
	// handler.

	var eve *github.PushEvent
	{
		eve = &github.PushEvent{
			Ref:   github.Ptr(branch + key.Bra),
			After: github.Ptr(one.Hash),
			Repo: &github.PushEventRepository{
				Name:  github.Ptr(key.Rep),
				Owner: &github.User{Login: github.Ptr(key.Org)},
			},
			HeadCommit: &github.HeadCommit{
				Timestamp: &github.Timestamp{Time: one.Time},
			},
		}
	}

	{
		err := whk.Update(context.Background(), "", "", eve)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Ensure that Webhook.Update stores our commit object in the cache.

	var sto Commit
	{
		sto, exi = whk.cac[key]
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	if !sto.Equals(one) {
		t.Fatal("expected", one.Hash, "got", sto.Hash)
	}

	// Once we stored the first commit in our internal cache, this first commit
	// object should still be considered the latest.

	{
		lat = whk.Latest(key, one)
	}

	if !one.Equals(lat) {
		t.Fatal("expected", one.Hash, "got", lat.Hash)
	}

	//
	// Second Commit 0x3456
	//

	var two Commit
	{
		two = Commit{
			Hash: "0x3456",
			Time: time.Unix(1, 0),
		}
	}

	{
		lat = whk.Latest(key, two)
	}

	// We have our first commit 0x1234 in the cache. When we ask for the latest by
	// comparing with an older commit, then we should get 0x1234 and not 0x3456.

	if two.Equals(lat) {
		t.Fatal("expected", one.Hash, "got", two.Hash)
	}

	// Calling Webhook.Update with the outdated commit should not change the
	// internal cache state.

	{
		eve = &github.PushEvent{
			Ref:   github.Ptr(branch + key.Bra),
			After: github.Ptr(two.Hash),
			Repo: &github.PushEventRepository{
				Name:  github.Ptr(key.Rep),
				Owner: &github.User{Login: github.Ptr(key.Org)},
			},
			HeadCommit: &github.HeadCommit{
				Timestamp: &github.Timestamp{Time: two.Time},
			},
		}
	}

	{
		err := whk.Update(context.Background(), "", "", eve)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		sto, exi = whk.cac[key]
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	if !sto.Equals(one) {
		t.Fatal("expected", one.Hash, "got", sto.Hash)
	}

	//
	// Third Commit 0x5678
	//

	var thr Commit
	{
		thr = Commit{
			Hash: "0x5678",
			Time: time.Unix(5, 0),
		}
	}

	{
		lat = whk.Latest(key, thr)
	}

	// We still have our first commit 0x1234 in the cache. When we ask for the
	// latest by comparing with a newer commit, then we should get 0x5678 and not
	// 0x1234.

	if !thr.Equals(lat) {
		t.Fatal("expected", thr.Hash, "got", lat.Hash)
	}

	// Calling Webhook.Update with the latest commit should change the internal
	// cache state.

	{
		eve = &github.PushEvent{
			Ref:   github.Ptr(branch + key.Bra),
			After: github.Ptr(thr.Hash),
			Repo: &github.PushEventRepository{
				Name:  github.Ptr(key.Rep),
				Owner: &github.User{Login: github.Ptr(key.Org)},
			},
			HeadCommit: &github.HeadCommit{
				Timestamp: &github.Timestamp{Time: thr.Time},
			},
		}
	}

	{
		err := whk.Update(context.Background(), "", "", eve)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		sto, exi = whk.cac[key]
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	if !sto.Equals(thr) {
		t.Fatal("expected", thr.Hash, "got", sto.Hash)
	}
}
