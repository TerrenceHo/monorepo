package stackerrors

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFuncName(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		input  string
		wanted string
	}

	testcases := []testcase{
		{
			input:  "github.com/org/path/package.FuncName",
			wanted: "FuncName",
		},
		{
			input:  "github.com/org/path/package.Class.FuncName",
			wanted: "Class.FuncName",
		},
		{
			input:  "github.com/org/path/package.(*Class).FuncName",
			wanted: "(*Class).FuncName",
		},
		{
			input:  "FuncName",
			wanted: "FuncName",
		},
		{
			input:  "",
			wanted: "",
		},
		{
			input:  "main.FakeCaller",
			wanted: "FakeCaller",
		},
	}

	for _, testcase := range testcases {
		s := processFuncName(testcase.input)
		assert.Equalf(
			s, testcase.wanted, "(processFuncName) %s != %s", s, testcase.wanted,
		)
	}
}

func TestFuncFromPC(t *testing.T) {
	assert := assert.New(t)

	wanted := []string{
		"helpTestFunc",
		"fakeStruct.helpTestMethod",
		"(*fakeStruct).helpTestPointerMethod",
	}

	fs := fakeStruct{}
	fsp := &fakeStruct{}

	hTF := helpTestFunc()
	assert.Equalf(wanted[0], hTF, "(TestFunc): %s != %s", wanted[0], hTF)
	hTM := fs.helpTestMethod()
	assert.Equalf(wanted[1], hTM, "(TestFunc): %s != %s", wanted[1], hTM)
	hTPM := fsp.helpTestPointerMethod()
	assert.Equalf(wanted[2], hTPM, "(TestFunc): %s != %s", wanted[2], hTPM)
}

func TestGetStackTrace(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		which    string
		function string
	}

	testcases := []testcase{
		{
			which:    "function",
			function: "fakeFunction",
		},
		{
			which:    "method",
			function: "stackTraceFake.fakeMethod",
		},
		{
			which:    "pointer-method",
			function: "(*stackTraceFake).fakePointerMethod",
		},
	}

	for _, testcase := range testcases {
		var st *stacktrace
		switch s := testcase.which; s {
		case "function":
			st = fakeFunction()
		case "method":
			st = stackTraceFake{}.fakeMethod()
		case "pointer-method":
			st = (&stackTraceFake{}).fakePointerMethod()
		default:
			assert.Fail("Testcase not found")
		}
		assert.Equalf(
			testcase.function, st.function,
			"(getStackTrace) %s != %s",
			testcase.function, st.function,
		)
		assert.Greaterf(st.line, 0, "(getStackTrace) %d <= %d", st.line, 0)
		assert.NotEmpty(st.file, "(getStackTrace) file was empty")
	}
}

func TestStacktraceFormat(t *testing.T) {
	assert := assert.New(t)
	type testcase struct {
		format string
		st     *stacktrace
		want   string
	}

	testcases := []testcase{
		{
			format: "%+v",
			st: &stacktrace{
				file:     "file/path.go",
				function: "funcName",
				line:     32,
			},
			want: "--- file/path.go:32 (funcName) ---",
		},
		{
			format: "string %+v more string",
			st: &stacktrace{
				file:     "file/path.go",
				function: "funcName",
				line:     32,
			},
			want: "string --- file/path.go:32 (funcName) --- more string",
		},
		{
			format: "string %v more string",
			st: &stacktrace{
				file:     "file/path.go",
				function: "funcName",
				line:     32,
			},
			want: "string  more string",
		},
	}

	for _, testcase := range testcases {
		s := fmt.Sprintf(testcase.format, testcase.st)
		assert.Equalf(testcase.want, s, "%s != %s", testcase.want, s)
	}
}

///// Helper test functions for TestFuncFromPC
func helpTestFunc() string {
	return funcFromPC(getPC())
}

type fakeStruct struct{}

func (f fakeStruct) helpTestMethod() string {
	return funcFromPC(getPC())
}

func (f *fakeStruct) helpTestPointerMethod() string {
	return funcFromPC(getPC())
}

func getPC() uintptr {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return 0
	}
	return pc
}

///// Helper test functions for TestGetStackTrace
func fakeErrorCaller() *stacktrace {
	return getStackTrace()
}

type stackTraceFake struct{}

func fakeFunction() *stacktrace {
	return fakeErrorCaller()
}

func (stf stackTraceFake) fakeMethod() *stacktrace {
	return fakeErrorCaller()
}

func (stf *stackTraceFake) fakePointerMethod() *stacktrace {
	return fakeErrorCaller()
}
