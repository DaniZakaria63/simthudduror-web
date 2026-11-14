package staff

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func newID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("generate id: %v", err))
	}
	return hex.EncodeToString(buf)
}
