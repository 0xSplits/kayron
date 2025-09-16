package preview

import (
	"bytes"
	"regexp"

	"github.com/xh3b4sd/tracer"
)

var (
	exp = regexp.MustCompile(`.*\/.*:(\$\{[^}]+\}|[A-Za-z0-9_][A-Za-z0-9._-]{0,127})($|["']?[ \t]*#.*|["'])`)
)

func repIma(lin []byte, tag []byte) ([]byte, error) {
	var sub [][]byte
	{
		sub = exp.FindSubmatch(lin)
		if len(sub) < 2 {
			return nil, tracer.Mask(containerImageFormatError, tracer.Context{Key: "input", Value: string(lin)})
		}
	}

	var out []byte
	{
		out = bytes.ReplaceAll(lin, sub[1], tag)
	}

	return out, nil
}
