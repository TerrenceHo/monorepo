package logging

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureStderr(t *testing.T, f func()) string {
	oldStderr := os.Stderr
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err.Error())
	}
	os.Stderr = write

	f()

	write.Close()
	data, err := ioutil.ReadAll(read)
	if err != nil {
		t.Fatal(err.Error())
	}

	os.Stderr = oldStderr
	return string(data)
}

func TestLoggers(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		loggerType    LoggerType
		hasMessages   bool
		hasCaller     bool
		hasStacktrace bool
		hasJSON       bool
	}

	testcases := []testcase{
		{
			loggerType:    TestLogger,
			hasMessages:   false,
			hasCaller:     false,
			hasStacktrace: false,
			hasJSON:       false,
		},
		{
			loggerType:    DevLogger,
			hasMessages:   true,
			hasCaller:     true,
			hasStacktrace: true,
			hasJSON:       false,
		},
		{
			loggerType:    ProdLogger,
			hasMessages:   true,
			hasCaller:     true,
			hasStacktrace: true,
			hasJSON:       true,
		},
	}

	for _, testcase := range testcases {
		output := captureStderr(t, func() {
			logger, err := ConfigureLogger(testcase.loggerType)
			assert.Nil(err)
			SetGlobalLogger(logger)
			Error("error message")
		})
		hasMessages := len(output) > 0
		assert.EqualValuesf(
			testcase.hasMessages, hasMessages,
			"loggerType %s hasMessages is wrong -- %s", testcase.loggerType, output,
		)
		hasCaller := strings.Contains(output, "testing.go")
		assert.EqualValuesf(
			testcase.hasCaller, hasCaller,
			"loggerType %s hasCaller is wrong -- %s", testcase.loggerType, output,
		)
		hasStacktrace := strings.Contains(output, "logging_test.go")
		assert.EqualValuesf(
			testcase.hasStacktrace, hasStacktrace,
			"loggerType %s hasStacktrace is wrong -- %s", testcase.loggerType, output,
		)
		hasJSON := strings.Contains(output, "{")
		assert.EqualValuesf(
			testcase.hasJSON, hasJSON,
			"loggerType %s hasJSON is wrong -- %s", testcase.loggerType, output,
		)
	}
}
