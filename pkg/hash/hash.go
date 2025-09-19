package hash

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Hash struct {
	// dsh is the dashed prefix version of Upp.
	//
	//     -1D0FD508
	//
	dsh []byte
	// low is the lower case version of Upp.
	//
	//     1d0fd508
	//
	low []byte
	// upp is the upper case version of this hash.
	//
	//     1D0FD508
	//
	upp []byte
}

func New(str string) Hash {
	var hsh string
	{
		hsh = newHsh(str)
	}

	var low string
	var upp string
	{
		low = cases.Lower(language.English).String(hsh)
		upp = cases.Upper(language.English).String(hsh)
	}

	return Hash{
		dsh: []byte("-" + upp),
		low: []byte(low),
		upp: []byte(upp),
	}
}

func (h Hash) Dashed() string {
	return string(h.dsh)
}

func (h Hash) Empty() bool {
	return h.dsh == nil && h.low == nil && h.upp == nil
}

func (h Hash) Lower() string {
	return string(h.low)
}

func (h Hash) Upper() string {
	return string(h.upp)
}

func newHsh(str string) string {
	sum := sha256.Sum256([]byte(str))
	enc := hex.EncodeToString(sum[:])

	return enc[:8]
}
