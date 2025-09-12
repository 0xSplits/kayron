package scanner

import (
	"bufio"
	"bytes"
)

func (s *Scanner) Append(pre []byte, suf []byte) *Scanner {
	var buf *bufio.Scanner
	{
		buf = bufio.NewScanner(bytes.NewReader(s.inp))
	}

	var blo [][]byte
	for buf.Scan() {
		var lin []byte
		{
			lin = append([]byte(nil), buf.Bytes()...) // copy to prevent buffer overwrites
		}

		if bytes.HasPrefix(lin, pre) {
			lin = insert(lin, suf)
		}

		{
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

func insert(lin []byte, suf []byte) []byte {
	var las byte
	{
		las = lin[len(lin)-1]
	}

	// double quote    0x22
	// single quote    0x27
	// colon           0x3A

	if las == 0x22 || las == 0x27 || las == 0x3A {
		return merge(lin[:len(lin)-1], suf, []byte{las})
	}

	return append(lin, suf...)
}

func merge(pre []byte, mid []byte, suf []byte) []byte {
	out := make([]byte, len(pre)+len(mid)+len(suf))

	num := copy(out, pre)
	num += copy(out[num:], mid)

	copy(out[num:], suf)

	return out
}
