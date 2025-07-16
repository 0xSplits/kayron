package roghfs

import (
	"os"
	"syscall"
)

// Chmod does nothing but returning the error syscall.EPERM.
func (r *Roghfs) Chmod(_ string, _ os.FileMode) error {
	return syscall.EPERM
}
