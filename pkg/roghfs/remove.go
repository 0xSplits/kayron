package roghfs

import "syscall"

// Remove does nothing but returning the error syscall.EPERM.
func (r *Roghfs) Remove(_ string) error {
	return syscall.EPERM
}
