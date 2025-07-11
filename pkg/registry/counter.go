package registry

func (r *Registry) Counter(nam string, val float64, lab map[string]string) error {
	return r.record(r.cou, nam, val, lab)
}
