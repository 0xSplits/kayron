package scanner

import (
	"bufio"
	"bytes"
	"unicode"
)

// Search tries to find the entire YAML block identified by the given key line,
// e.g. "  Service:". A new scanner configured with the found YAML block as
// input bytes is returned.
func (s *Scanner) Search(key []byte) *Scanner {
	var buf *bufio.Scanner
	{
		buf = bufio.NewScanner(bytes.NewReader(s.inp))
	}

	var blo [][]byte
	var fou bool
	var end int
	var sta int
	for buf.Scan() {
		var lin []byte
		{
			lin = append([]byte(nil), buf.Bytes()...) // copy to prevent buffer overwrites
		}

		if fou {
			end = spaces(lin)
		}

		if fou && sta == end && len(lin) != 0 {
			break
		}

		if !fou {
			fou = bytes.HasPrefix(lin, key)
			sta = spaces(lin)
		}

		if fou {
			blo = append(blo, lin)
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

func spaces(b []byte) int {
	var cou int

	for _, x := range b {
		if unicode.IsSpace(rune(x)) {
			cou++
		} else {
			break
		}
	}

	return cou
}
