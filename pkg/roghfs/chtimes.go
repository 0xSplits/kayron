package roghfs

import (
	"syscall"
	"time"
)

// Chtimes does nothing but returning the error syscall.EPERM.
func (r *Roghfs) Chtimes(n string, _ time.Time, _ time.Time) error {
	return syscall.EPERM
}
