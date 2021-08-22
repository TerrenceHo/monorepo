package hash

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/TerrenceHo/monorepo/utils-go/random"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = stackerrors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = stackerrors.New("incompatible version of argon2")
)

// Argon2ID implements the Hasher interface. Underneath, it uses the Argon2ID
// implementation to hash passwords and compare passwords.
type Argon2ID struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
	version     int
}

func NewArgon2ID(memory, iterations uint32, parallelism uint8, saltLength, keyLength uint32) *Argon2ID {
	return &Argon2ID{
		memory:      memory,
		iterations:  iterations,
		parallelism: parallelism,
		saltLength:  saltLength,
		keyLength:   keyLength,
		version:     argon2.Version,
	}
}

// HashPassword takes in a password string and hashes + salts it using the
// Argon2ID algorithm.  It is recommended that you pass in a password + pepper
// combination, because the pepper should be a secret value that is not stored
// alongside the hashed password.
//
// The password is stored in the following format:
//		$argon2id$v={version}$m={memory},t={iterations},p={parallelism}${salt}${hash}
func (a *Argon2ID) HashPassword(password string) (encodedHash string, err error) {
	// Generate a cryptographically secure random salt.
	salt, err := random.GenerateRandomBytes(a.saltLength)
	if err != nil {
		return "", err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey([]byte(password), salt, a.iterations, a.memory, a.parallelism, a.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, a.memory, a.iterations, a.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// ComparePasswordAndHash takes a password and the hash, and extracts the salt
// and strength parameters, and compares the hashes in constant time (avoiding
// time collision attacks).
func (a *Argon2ID) ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	salt, hash, err := a.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, a.iterations, a.memory, a.parallelism, a.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (a *Argon2ID) decodeHash(encodedHash string) (salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return nil, nil, err
	}
	if memory != a.memory || iterations != a.iterations || parallelism != a.parallelism {
		return nil, nil, ErrInvalidHash
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}
	saltLength := uint32(len(salt))
	if saltLength != a.saltLength {
		return nil, nil, ErrInvalidHash
	}

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
	}
	keyLength := uint32(len(hash))
	if keyLength != a.keyLength {
		return nil, nil, ErrInvalidHash
	}

	return salt, hash, nil
}
