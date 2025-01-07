package str

import (
	"crypto/rand"
	"fmt"
)

func Random() string {
	b := make([]byte, 3)
	rand.Read(b)

	return fmt.Sprintf("%x", b)
}