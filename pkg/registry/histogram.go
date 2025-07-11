package registry

func (r *Registry) Histogram(nam string, val float64, lab map[string]string) error {
	return r.record(r.his, nam, val, lab)
}
