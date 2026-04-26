package hash

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
)

func TestArgon2IDHashPassword(t *testing.T) {
	assert := assert.New(t)
	password := "password" // the classical bad password
	var memory uint32 = 64 * 1024
	var iterations uint32 = 3
	var parallelism uint8 = 4
	var saltLength uint32 = 16
	var keyLength uint32 = 32

	a := NewArgon2ID(
		memory, iterations, parallelism, saltLength, keyLength,
	)

	hash1, err := a.HashPassword(password)
	hash2, err := a.HashPassword(password)
	assert.Nil(err)
	assert.NotEqualValues(hash1, password, "password should not equal hash1")
	assert.NotEqualValues(hash2, password, "password should not equal hash2")
	assert.NotEqualValues(hash1, hash2, "hash1 should not equal hash2")

	decomposed := strings.Split(hash1, "$")
	assert.Len(decomposed, 6, "decomposed hash should have six parts")
	assert.EqualValues("", decomposed[0])
	assert.EqualValues("argon2id", decomposed[1], "decomposed hash name not correct")
	assert.EqualValues("v="+strconv.Itoa(argon2.Version), decomposed[2], "decomposed hash version not equal")
	assert.EqualValuesf(
		fmt.Sprintf("m=%d,t=%d,p=%d", memory, iterations, parallelism),
		decomposed[3],
		"%s != %s",
		fmt.Sprintf("m=%d,t=%d,p=%d", memory, iterations, parallelism),
		decomposed[3],
	)
	assert.NotEqualValues(password, decomposed[5], "decomposed hash password and hash should not be equal")
}

func getHash(
	t *testing.T,
	hashFunc func(password string) (string, error),
	password string,
) string {
	assert := assert.New(t)
	encodedHash, err := hashFunc(password)
	assert.Nil(err)
	return encodedHash
}

func TestArgon2IDComparePasswordAndHash(t *testing.T) {
	assert := assert.New(t)
	var memory uint32 = 64 * 1024
	var iterations uint32 = 3
	var parallelism uint8 = 4
	var saltLength uint32 = 16
	var keyLength uint32 = 32
	a := NewArgon2ID(
		memory, iterations, parallelism, saltLength, keyLength,
	)

	type testcase struct {
		password string
		hash     string
		compare  bool
		err      error
	}

	hash := getHash(t, a.HashPassword, "password")

	testcases := []testcase{
		{
			password: "password",
			hash:     hash,
			compare:  true,
			err:      nil,
		},
		{
			password: "password",
			hash:     "fake hash",
			compare:  false,
			err:      ErrInvalidHash,
		},
		{
			password: "password",
			hash:     strings.Replace(hash, "argon2id", "argon3id", 1),
			compare:  false,
			err:      ErrWrongHashAlgorithm,
		},
		{
			password: "password",
			hash:     strings.Replace(hash, "v="+strconv.Itoa(argon2.Version), "v=1", 1),
			compare:  false,
			err:      ErrIncompatibleVersion,
		},
		{
			password: "password",
			hash:     strings.Replace(hash, "m="+strconv.Itoa(int(memory)), "m=1", 1),
			compare:  false,
			err:      ErrInvalidHash,
		},
		{
			password: "password",
			hash:     strings.Replace(hash, "t="+strconv.Itoa(int(iterations)), "t=1", 1),
			compare:  false,
			err:      ErrInvalidHash,
		},
		{
			password: "password",
			hash:     strings.Replace(hash, "p="+strconv.Itoa(int(parallelism)), "p=1", 1),
			compare:  false,
			err:      ErrInvalidHash,
		},
		{
			password: "password",
			hash:     hash[:len(hash)-int(keyLength)-2] + hash[len(hash)-int(keyLength)-1:],
			compare:  false,
			err:      ErrInvalidHash,
		},
		{
			password: "password",
			hash:     hash[:len(hash)-1],
			compare:  false,
			err:      ErrInvalidHash,
		},
	}

	for _, testcase := range testcases {
		compare, err := a.ComparePasswordAndHash(testcase.password, testcase.hash)
		assert.EqualValuesf(
			testcase.compare,
			compare,
			"hash comparison not correct",
		)
		if testcase.err != nil {
			assert.EqualError(testcase.err, err.Error())
		} else {
			assert.Nil(err)
		}
	}

}
