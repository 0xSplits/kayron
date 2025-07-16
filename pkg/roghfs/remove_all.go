package roghfs

import "syscall"

// RemoveAll does nothing but returning the error syscall.EPERM.
func (r *Roghfs) RemoveAll(_ string) error {
	return syscall.EPERM
}
