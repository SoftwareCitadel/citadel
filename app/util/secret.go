package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

// GenerateSecretKey generates a random secret key of the specified length.
func GenerateSecretKey() (string, error) {
	const (
		MIN_LENGTH = 32
		MAX_LENGTH = 64
	)

	// Generate a random number between min and max.
	randomBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(MAX_LENGTH-MIN_LENGTH+1)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random length: %v", err)
	}

	// Convert to int and adjust for the range offset.
	length := int(randomBigInt.Int64()) + MIN_LENGTH

	if length <= 0 {
		return "", fmt.Errorf("length must be a positive integer")
	}

	// Create a byte slice to hold the random bytes.
	key := make([]byte, length)

	// Read random bytes from the crypto/rand reader.
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate secret key: %v", err)
	}

	// Convert the byte slice to a hexadecimal string.
	secretKey := hex.EncodeToString(key)
	return secretKey, nil
}
