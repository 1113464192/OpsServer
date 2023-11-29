package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"strings"
)

func Validate(h func() hash.Hash, message, key []byte, signature string) bool {
	decoded, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return SecondValidate(h, message, key, decoded)
}

func ValidatePrefix(message, key []byte, signature string) bool {
	parts := strings.Split(signature, "=")
	if len(parts) != 2 {
		return false
	}
	switch parts[0] {
	case "sha1":
		return Validate(sha1.New, message, key, parts[1])
	case "sha256":
		return Validate(sha256.New, message, key, parts[1])
	default:
		return false
	}
}

func SecondValidate(h func() hash.Hash, message, key, signature []byte) bool {
	mac := hmac.New(h, key)
	mac.Write(message)
	sum := mac.Sum(nil)
	return hmac.Equal(signature, sum)
}
