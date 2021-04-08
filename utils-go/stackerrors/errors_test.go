package stackerrors

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		err_ctx string
	}

	testcases := []testcase{
		{err_ctx: "new error context"},
		{err_ctx: "second error context"},
	}

	for _, testcase := range testcases {
		e := New(testcase.err_ctx)
		assert.EqualErrorf(e, testcase.err_ctx, "New error %w != %s", e, testcase.err_ctx)
	}
}

func TestErrorf(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		format string
		args   []interface{}
		want   string
	}

	testcases := []testcase{
		{
			format: "%s %s",
			args: []interface{}{
				"error", "string",
			},
			want: "error string",
		},
		{
			format: "%v errors in %s",
			args: []interface{}{
				5, "package",
			},
			want: "5 errors in package",
		},
	}

	for _, testcase := range testcases {
		e := Errorf(testcase.format, testcase.args...)
		assert.EqualErrorf(
			e, testcase.want, "Errorf error %v != %s", e, testcase.want,
		)
	}
}

func TestWrap(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		err    error
		format string
		want   string
	}

	testcases := []testcase{
		{
			err:    New("base error"),
			format: "read error",
			want:   "read error: base error",
		},
		{
			err:    Wrap(New("base error"), "second error"),
			format: "third error",
			want:   "third error: second error: base error",
		},
		{
			err:    errors.New("std error"),
			format: "wrapper error",
			want:   "wrapper error: std error",
		},
	}

	for _, testcase := range testcases {
		e := Wrap(testcase.err, testcase.format)
		assert.EqualErrorf(
			e, testcase.want, "Wrap error %v != %s", e, testcase.want,
		)
	}
}

func TestWrapf(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		err    error
		format string
		args   []interface{}
		want   string
	}

	testcases := []testcase{
		{
			err:    New("base error"),
			format: "%s number %d",
			args: []interface{}{
				"error", 2,
			},
			want: "error number 2: base error",
		},
		{
			err:    Wrap(New("base error"), "second error"),
			format: "error layer %d",
			args: []interface{}{
				3,
			},
			want: "error layer 3: second error: base error",
		},
		{
			err:    errors.New("std error"),
			format: "wrapper error",
			args:   []interface{}{},
			want:   "wrapper error: std error",
		},
	}

	for _, testcase := range testcases {
		e := Wrapf(testcase.err, testcase.format, testcase.args...)
		assert.EqualErrorf(
			e, testcase.want, "Wrap error %v != %s", e, testcase.want,
		)
	}
}

func TestIs(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		err    error
		target error
		want   bool
	}

	testcases := []testcase{
		{
			err:    New("error"),
			target: New("error"),
			want:   true,
		},
		{
			err:    New("error"),
			target: errors.New("error"),
			want:   true,
		},
		{
			err:    Wrap(New("error"), "second error"),
			target: Wrap(New("error"), "second error"),
			want:   true,
		},
		{
			err:    Wrap(New("error"), "second error"),
			target: New("error"),
			want:   true,
		},
		{
			err:    New("error"),
			target: Wrap(New("error"), "second error"),
			want:   false,
		},
		{
			err:    Wrap(Wrap(New("error"), "error"), "error"),
			target: New("error"),
			want:   true,
		},
		{
			err:    Wrap(Wrap(New("error"), "error"), "error"),
			target: Wrap(New("error"), "error"),
			want:   true,
		},
	}

	for _, testcase := range testcases {
		isEqual := errors.Is(testcase.err, testcase.target)

		assert.Equalf(
			isEqual,
			testcase.want,
			"errors.Is did not succeed: %b != %b for error %v and error %v",
			isEqual,
			testcase.want,
			testcase.err,
			testcase.target,
		)
	}
}

func TestAs(t *testing.T) {
	assert := assert.New(t)
	type testcase struct {
		input     error
		want_bool bool
		want_err  error
	}

	testcases := []testcase{
		{
			input:     New("base error"),
			want_bool: true,
			want_err:  New("base error"),
		},
		{
			input:     Wrapf(New("base error"), "second layer"),
			want_bool: true,
			want_err:  Wrapf(New("base error"), "second layer"),
		},
		{
			input:     fmt.Errorf("fmt base error"),
			want_bool: false,
			want_err:  nil,
		},
	}

	var ce *contextError
	for _, testcase := range testcases {
		b := errors.As(testcase.input, &ce)
		if b {
			assert.EqualErrorf(
				ce, testcase.want_err.Error(),
				"(As) %s != %s", ce, testcase.want_err.Error(),
			)
		} else {
			assert.False(testcase.want_bool)
		}
	}
}

func TestUnwrap(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		input error
		want  error
	}

	testcases := []testcase{
		{
			input: nil,
			want:  nil,
		},
		{
			input: New("no unwrap error"),
			want:  nil,
		},
		{
			input: Wrap(New("base error"), "second layer"),
			want:  New("base error"),
		},
		{
			input: Wrap(Wrap(New("base error"), "second layer"), "third layer"),
			want:  Wrap(New("base error"), "second layer"),
		},
	}

	for _, testcase := range testcases {
		e := errors.Unwrap(testcase.input)
		if e == nil {
			assert.Nil(testcase.want)
		} else {
			assert.EqualErrorf(
				e,
				testcase.want.Error(),
				"TestUnwrap: %s != %s",
				testcase.input.Error(),
				testcase.want.Error(),
			)
		}
	}
}

func TestRootCause(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		input error
		want  error
	}

	testcases := []testcase{
		{
			input: nil,
			want:  nil,
		},
		{
			input: New("base error"),
			want:  New("base error"),
		},
		{
			input: fmt.Errorf("new error"),
			want:  fmt.Errorf("new error"),
		},
		{
			input: Wrap(Wrap(New("base error"), "level 2"), "level 3"),
			want:  New("base error"),
		},
		{
			input: Wrap(&contextError{
				ctx: "context",
				err: nil,
			}, "second layer"),
			want: New("context"),
		},
	}

	for _, testcase := range testcases {
		e := RootCause(testcase.input)
		if e == nil {
			assert.Nil(testcase.want)
		} else {
			assert.EqualErrorf(
				e, testcase.want.Error(),
				"(RootCause) %s != %s",
				e.Error(), testcase.want.Error(),
			)
		}
	}
}

func TestContextFormat(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		format string
		err    error
		want   string
	}

	testcases := []testcase{
		{
			format: "%v",
			err:    New("base error"),
			want:   "base error",
		},
		{
			format: "%s",
			err:    New("base error"),
			want:   "base error",
		},
		{
			format: "%q",
			err:    New("base error"),
			want:   "\"base error\"",
		},
		{
			format: "%v",
			err:    Wrap(New("base error"), "second layer"),
			want:   "second layer: base error",
		},
		{
			format: "%s",
			err:    Wrap(New("base error"), "second layer"),
			want:   "second layer: base error",
		},
		{
			format: "%q",
			err:    Wrap(New("base error"), "second layer"),
			want:   "\"second layer: base error\"",
		},
		{
			format: "%v",
			err:    Wrap(Wrap(New("base error"), "second layer"), "third layer"),
			want:   "third layer: second layer: base error",
		},
		{
			format: "%s",
			err:    Wrap(Wrap(New("base error"), "second layer"), "third layer"),
			want:   "third layer: second layer: base error",
		},
		{
			format: "%q",
			err:    Wrap(Wrap(New("base error"), "second layer"), "third layer"),
			want:   "\"third layer: second layer: base error\"",
		},
	}

	for _, testcase := range testcases {
		s := fmt.Sprintf(testcase.format, testcase.err)
		assert.Equalf(
			testcase.want, s, "(ContextFormat) %s != %s", testcase.want,
		)
	}
}

func TestFormatStackTrace(t *testing.T) {
	assert := assert.New(t)

	type wantedTrace struct {
		msg      string
		file     string
		function string
	}
	type testcase struct {
		format string
		err    error
		want   []wantedTrace
	}

	err1 := New("base error")
	err2 := Wrap(err1, "second layer")
	err3 := Wrap(err2, "third layer")

	testcases := []testcase{
		{
			format: "%+v",
			err:    err1,
			want: []wantedTrace{
				{
					msg:      "base error",
					file:     "utils-go/stackerrors/errors_test.go",
					function: "TestFormatStackTrace",
				},
			},
		},
		{
			format: "%+v",
			err:    err2,
			want: []wantedTrace{
				{
					msg:      "second layer",
					file:     "utils-go/stackerrors/errors_test.go",
					function: "TestFormatStackTrace",
				},
				{
					msg:      "base error",
					file:     "utils-go/stackerrors/errors_test.go",
					function: "TestFormatStackTrace",
				},
			},
		},
		{
			format: "%+v",
			err:    err3,
			want: []wantedTrace{
				{
					msg:      "third layer",
					file:     "utils-go/stackerrors/errors_test.go",
					function: "TestFormatStackTrace",
				},
				{
					msg:      "second layer",
					file:     "utils-go/stackerrors/errors_test.go",
					function: "TestFormatStackTrace",
				},
				{
					msg:      "base error",
					file:     "utils-go/stackerrors/errors_test.go",
					function: "TestFormatStackTrace",
				},
			},
		},
	}

	for _, testcase := range testcases {
		s := fmt.Sprintf(testcase.format, testcase.err)
		t.Log(s)
		msgs := strings.Split(s, "\n")
		for i := range testcase.want {
			wt := testcase.want[i]
			assert.Equal(wt.msg, msgs[0])
			assert.Truef(
				strings.Contains(msgs[1], wt.file),
				"(file) %s does not contain %s",
				msgs[1], wt.file,
			)
			assert.Truef(
				strings.Contains(msgs[1], wt.function),
				"(file) %s does not contain %s",
				msgs[1], wt.file,
			)
			msgs = msgs[2:]
		}
	}
}
