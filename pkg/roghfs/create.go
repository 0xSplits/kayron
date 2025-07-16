package roghfs

import (
	"syscall"

	"github.com/spf13/afero"
)

// Create does nothing but returning the error syscall.EPERM.
func (r *Roghfs) Create(_ string) (afero.File, error) {
	return nil, syscall.EPERM
}
