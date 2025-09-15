package scanner

import (
	"bufio"
	"bytes"
)

// Delete tries to drop the entire YAML block identified by the given key line,
// e.g. "  Service:". A new scanner configured without the found YAML block as
// input bytes is returned. Optionally a substituion slice may be defined, which
// causes Delete to replace the matched block with the given structure.
func (s *Scanner) Delete(key []byte, sub ...byte) *Scanner {
	var buf *bufio.Scanner
	{
		buf = bufio.NewScanner(bytes.NewReader(s.inp))
	}

	var blo [][]byte
	var drp bool
	var rep bool
	var end int
	var sta int
	for buf.Scan() {
		var lin []byte
		{
			lin = append([]byte(nil), buf.Bytes()...) // copy to prevent buffer overwrites
		}

		if drp {
			end = spaces(lin)
		}

		if drp && end <= sta && len(lin) != 0 {
			drp = false
		}

		if bytes.HasPrefix(lin, key) {
			drp = true
			sta = spaces(lin)
		}

		if !drp {
			blo = append(blo, lin)
		} else if !rep && sub != nil {
			rep = true
			blo = append(blo, sub)
		}
	}

	var inp []byte
	{
		inp = bytes.Join(blo, []byte("\n"))
	}

	return New(Config{
		Inp: inp,
	})
}
