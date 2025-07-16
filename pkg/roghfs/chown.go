package roghfs

import "syscall"

// Chown does nothing but returning the error syscall.EPERM.
func (r *Roghfs) Chown(_ string, _ int, _ int) error {
	return syscall.EPERM
}
