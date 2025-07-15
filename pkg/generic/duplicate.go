package generic

// Duplicate returns a list of the duplicated items within the given list.
func Duplicate[T comparable](lis []T) []T {
	all := map[T]struct{}{}
	tmp := map[T]struct{}{}

	for _, x := range lis {
		{
			_, exi := all[x]
			if exi {
				tmp[x] = struct{}{}
			}
		}

		{
			all[x] = struct{}{}
		}
	}

	var dup []T
	for k := range tmp {
		dup = append(dup, k)
	}

	return dup
}
