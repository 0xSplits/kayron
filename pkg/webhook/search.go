package webhook

import "path"

func (w *Webhook) Search(own string, rep string, bra string, cmp Commit) Commit {
	// "owner/repo/branch"
	var key string
	{
		key = path.Join(own, rep, bra)
	}

	// be explicit about empty caches

	var cac Commit
	var exi bool
	{
		cac, exi = w.search(key)
		if !exi {
			return cmp
		}
	}

	// if the cached commit is newer, return the cached version
	if cac.Time.After(cmp.Time) {
		return cac
	}

	// if the cached commit is older or equally old, return the compare version

	return cmp
}

func (w *Webhook) search(key string) (Commit, bool) {
	w.mut.Lock()
	cac, exi := w.cac[key]
	w.mut.Unlock()
	return cac, exi
}
