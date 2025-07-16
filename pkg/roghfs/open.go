package roghfs

import (
	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

// Open tries to open the given file. Note that Open() is called after every
// walk function loop for directories when using afero.Walk().
func (r *Roghfs) Open(pat string) (afero.File, error) {
	var err error

	{
		err = r.ensure(pat)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var fil afero.File
	{
		fil, err = r.bas.Open(pat)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return fil, nil
}
