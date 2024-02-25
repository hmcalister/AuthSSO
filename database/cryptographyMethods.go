package database

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen    uint32 = 32
	timeCost   uint32 = 1
	memoryCost uint32 = 8 * 1024
	threads    uint8  = 1
	keyLen     uint32 = 32
)

// Generate a new salt and return it.
//
// This function can error if a kernel function errors, although this should never happen.
func generateSalt() (string, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	return string(salt), nil
}

// Perform the hash of a (plaintext) password with salt.
func calculateHash(password string, salt string) string {
	hash := argon2.IDKey([]byte(password), []byte(salt), timeCost, memoryCost, threads, keyLen)

	return string(hash)
}
