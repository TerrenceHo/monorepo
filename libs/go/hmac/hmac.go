package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type HMAC struct {
	hmac hash.Hash
}

// NewHMAC returns a new HMAC object
func NewHMAC(key string) HMAC {
	return HMAC{
		hmac: hmac.New(sha256.New, []byte(key)),
	}
}

// Hash uses the internal hash to compute a sum and returns the base64 encoded
// string.
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
