package hash

import (
	"crypto/sha256"
	"fmt"
)

func CalculateHash(hString string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hString)))
}
