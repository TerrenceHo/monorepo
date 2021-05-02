package path

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	assert := assert.New(t)
	switch runtime.GOOS {
	case "windows":
		assert.EqualValuesf(
			os.Getenv("USERPROFILE"), Home(),
			"%s != %s", os.Getenv("USERPROFILE"), Home(),
		)
	default:
		assert.EqualValuesf(
			os.Getenv("HOME"), Home(),
			"%s != %s", os.Getenv("HOME"), Home(),
		)
	}
}

func TestExpandUser(t *testing.T) {
	assert := assert.New(t)
	type testcase struct {
		input string
		want  string
	}
	t.Log("HOME: " + os.Getenv("HOME"))

	testcases := []testcase{
		{
			input: "~",
			want:  os.Getenv("HOME"),
		},
		{
			input: "~/go",
			want:  filepath.Join(os.Getenv("HOME"), "go"),
		},
		{
			input: "/fake~/go",
			want:  "/fake~/go",
		},
		{
			input: "/fake/go",
			want:  "/fake/go",
		},
	}

	for _, testcase := range testcases {
		assert.EqualValuesf(
			testcase.want, ExpandUser(testcase.input),
			"%s != %s", testcase.want, ExpandUser(testcase.input),
		)
	}
}
