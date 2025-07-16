package roghfs

import (
	"os"
	"syscall"

	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

// OpenFile opens a file using the given flags and the given permissions. The
// error syscall.EPERM is returned if the provided flags request any form of
// write access.
func (r *Roghfs) OpenFile(pat string, flg int, prm os.FileMode) (afero.File, error) {
	var err error

	if flg&(os.O_WRONLY|syscall.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 {
		return nil, syscall.EPERM
	}

	{
		err = r.ensure(pat)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var fil afero.File
	{
		fil, err = r.bas.OpenFile(pat, flg, prm)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return fil, nil
}
