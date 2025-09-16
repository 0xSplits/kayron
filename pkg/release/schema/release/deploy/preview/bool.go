package preview

import "strconv"

// Bool enables preview deployments for this service.
type Bool bool

func (b Bool) Empty() bool {
	return !bool(b)
}

func (b Bool) String() string {
	return strconv.FormatBool(bool(b))
}

func (b Bool) Verify() error {
	return nil
}
