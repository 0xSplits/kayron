package generic

// Duplicate expresses whether the given list contains duplicated items.
func Duplicate[T comparable](lis []T) bool {
	see := map[T]struct{}{}

	for _, x := range lis {
		{
			_, exi := see[x]
			if exi {
				return true
			}
		}

		{
			see[x] = struct{}{}
		}
	}

	return false
}
