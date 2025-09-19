package suspend

import "strconv"

// Bool disables any further reconciliation of this release indefinitely.
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
