package auth

import (
	"crypto/rand"
	"encoding/base32"
)

// GenerateToken generates 32-bit token
func GenerateToken() string {
	randomBytes := make([]byte, 32)
	_, _ = rand.Read(randomBytes)
	return base32.StdEncoding.EncodeToString(randomBytes)
}
