package registry

func (r *Registry) Gauge(nam string, val float64, lab map[string]string) error {
	return r.record(r.gau, nam, val, lab)
}
