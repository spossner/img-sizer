package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// calculateETag generates an ETag based on the image content and processing parameters
func CalculateETag(img []byte, params string) string {
	// Calculate SHA-256 hash of the image content and parameters
	hash := sha256.New()
	hash.Write(img)
	hash.Write([]byte(params))

	// Return the hash as a hex string
	return fmt.Sprintf("\"%s\"", hex.EncodeToString(hash.Sum(nil)))
}
