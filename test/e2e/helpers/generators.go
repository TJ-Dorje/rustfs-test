package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func RandomBucketName() string {
	return fmt.Sprintf("test-bucket-%s", randomHex(4))
}

func RandomKey() string {
	return fmt.Sprintf("test-object-%s", randomHex(4))
}

func RandomPayload(sizeBytes int) []byte {
	b := make([]byte, sizeBytes)
	_, _ = rand.Read(b)
	return b
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
