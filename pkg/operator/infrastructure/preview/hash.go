package preview

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(str string) string {
	sum := sha256.Sum256([]byte(str))
	enc := hex.EncodeToString(sum[:])
	return enc[:8]
}
