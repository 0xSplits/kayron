package handler

// Names returns the list of worker handler names matching the given list of
// worker handlers. E.g. this function should enable the worker's metrics
// registry to whitelist all worker handlers in the pkg/worker/handler package.
func Names(han []Interface) []string {
	var lis []string

	for _, x := range han {
		lis = append(lis, Name(x))
	}

	return lis
}
