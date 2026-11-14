package account

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const credentialBytes = 32

// NewCredential returns a cryptographically random opaque bearer/session
// credential.
func NewCredential() (string, error) {
	buf := make([]byte, credentialBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate credential: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
