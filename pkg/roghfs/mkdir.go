package roghfs

import (
	"os"
	"syscall"
)

// Mkdir does nothing but returning the error syscall.EPERM.
func (r *Roghfs) Mkdir(_ string, _ os.FileMode) error {
	return syscall.EPERM
}
