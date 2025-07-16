package roghfs

import "syscall"

// Rename does nothing but returning the error syscall.EPERM.
func (r *Roghfs) Rename(_ string, _ string) error {
	return syscall.EPERM
}
