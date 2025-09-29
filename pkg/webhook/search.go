package webhook

import "path"

func (w *Webhook) Search(own string, rep string, bra string, cmp Commit) Commit {
	// "owner/repo/branch"
	key := path.Join(own, rep, bra)

	w.mut.Lock()

	// be explicit about empty caches

	var cac Commit
	var exi bool
	{
		cac, exi = w.cac[key]
		if !exi {
			return cmp
		}
	}

	w.mut.Unlock()

	// if the cached commit is newer, return the cached version
	if cac.Time.After(cmp.Time) {
		return cac
	}

	// if the cached commit is older or equally old, return the compare version

	return cmp
}
