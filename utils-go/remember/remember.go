package remember

import (
	"encoding/base64"

	"github.com/TerrenceHo/monorepo/utils-go/random"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
)

const RememberTokenByteCount = 32

// RememberToken takes in a byteCount, and returns a random base64 encoded
// string of bytes. Recommended to use at least 32 byteCount.
func RememberToken(byteCount uint32) (string, error) {
	b, err := random.GenerateRandomBytes(byteCount)
	if err != nil {
		return "", stackerrors.Wrap(err, "generating remember token failed")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
