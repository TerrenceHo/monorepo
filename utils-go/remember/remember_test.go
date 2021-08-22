package remember

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemember(t *testing.T) {
	assert := assert.New(t)
	token, err := RememberToken(RememberTokenByteCount)
	assert.Nil(err, "Generating remember tokens should not error")
	count, err := countBytes(token)
	assert.Nil(err, "token is not base64 encoded")
	assert.EqualValues(
		RememberTokenByteCount,
		count,
		"remember token count length not equal to expected length",
	)
}

func countBytes(b64str string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(b64str)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}
