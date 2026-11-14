package post

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"
)

func trimSpace(value string) string { return strings.TrimSpace(value) }

func containsControl(value string) bool {
	return strings.ContainsFunc(value, func(r rune) bool {
		return unicode.IsControl(r) && r != '\n' && r != '\t'
	})
}

func newID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("generate id: %v", err))
	}
	return hex.EncodeToString(buf)
}
