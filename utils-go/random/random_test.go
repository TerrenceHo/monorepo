package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomBytes(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		size uint32
	}

	testcases := []TestCase{
		{
			size: 2,
		},
		{
			size: 4,
		},
		{
			size: 20,
		},
	}

	for _, testcase := range testcases {
		bytes, err := GenerateRandomBytes(testcase.size)
		assert.Nil(err)
		assert.Equalf(
			testcase.size,
			uint32(len(bytes)),
			"size of random bytes not equal: %d != %d",
			testcase.size,
			uint32(len(bytes)),
		)
	}
}
