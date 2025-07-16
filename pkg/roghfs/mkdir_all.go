package roghfs

import (
	"os"
	"syscall"
)

// MkdirAll does nothing but returning the error syscall.EPERM.
func (r *Roghfs) MkdirAll(_ string, _ os.FileMode) error {
	return syscall.EPERM
}
