package webhook

// Latest returns the most recent commit within a branch by comparing the
// already cached version, if any, with the compared version cmp. If the cache
// is empty for the underlying branch key, then cmp is returned.
func (w *Webhook) Latest(key Key, cmp Commit) Commit {
	// Be explicit about empty caches and just return the compared version if the
	// given key holds no value at this point.

	var cac Commit
	var exi bool
	{
		cac, exi = w.search(key)
		if !exi {
			return cmp
		}
	}

	// If the cached commit is newer, then return the cached version.

	if cac.Time.After(cmp.Time) {
		return cac
	}

	// If the cached commit is older or equally old, then return the compared
	// version.

	return cmp
}

func (w *Webhook) search(key Key) (Commit, bool) {
	w.mut.Lock()
	cac, exi := w.cac[key]
	w.mut.Unlock()
	return cac, exi
}
