package policy

// Delete purges the underlying local cache, causing Cancel to fetch the latest
// version of the stack object state again over network.
func (p *Policy) Delete() {
	p.mut.Lock()
	p.sta = nil
	p.mut.Unlock()
}
