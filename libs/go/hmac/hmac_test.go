package hmac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	hashInput  = "input"
	wrongInput = "wrong"
)

func TestHMAC(t *testing.T) {
	assert := assert.New(t)
	h := NewHMAC("key")
	hash := h.Hash(hashInput)
	wrongHash := h.Hash(wrongInput)
	secondHash := h.Hash(hashInput)
	assert.NotEqualValues(hash, wrongHash, "hashes should not have equal values")
	assert.EqualValues(hash, secondHash, "hashes should have equal values")
}
