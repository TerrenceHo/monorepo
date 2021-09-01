package fastlinks

import (
	"testing"

	"github.com/TerrenceHo/monorepo/utils-go/logging"
	"github.com/stretchr/testify/assert"
)

func TestLoggerType(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		input  string
		output logging.LoggerType
	}
	testcases := []testcase{
		{
			input:  "dev",
			output: logging.DevLogger,
		},
		{
			input:  "prod",
			output: logging.ProdLogger,
		},
		{
			input:  "testing",
			output: logging.TestLogger,
		},
		{
			input:  "fake",
			output: logging.TestLogger,
		},
	}

	for _, testcase := range testcases {
		assert.EqualValuesf(
			testcase.output, loggerType(testcase.input),
			"%s != %s", testcase.output, loggerType(testcase.input),
		)
	}
}
