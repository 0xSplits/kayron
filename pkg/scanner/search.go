package scanner

import (
	"bufio"
	"bytes"
	"unicode"
)

// Search tries to find the entire YAML block identified by the given key line,
// e.g. "  Service:".
func (s *Scanner) Search(key []byte) []byte {
	var sca *bufio.Scanner
	{
		sca = bufio.NewScanner(bytes.NewReader(s.inp))
	}

	var blo [][]byte
	var fou bool
	var end int
	var sta int
	for sca.Scan() {
		var lin []byte
		{
			lin = sca.Bytes()
		}

		if fou {
			end = spaces(lin)
		}

		if fou && sta == end && len(lin) != 0 {
			break
		}

		if !fou {
			fou = bytes.Equal(lin, key)
			sta = spaces(lin)
		}

		if fou {
			blo = append(blo, lin)
		}
	}

	var res []byte
	{
		res = bytes.Join(blo, []byte("\n"))
	}

	return res
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
