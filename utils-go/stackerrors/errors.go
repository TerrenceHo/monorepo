package stackerrors

import (
	"fmt"
	"io"
	"strings"
)

// New creates a new error with a context and a stack trace. Returns a
// contextError pointer.
func New(ctx string) error {
	return &contextError{
		err:        nil,
		ctx:        ctx,
		stacktrace: getStackTrace(),
	}
}

// Errorf is a convenience function mimics fmt.Errorf, but returns a
// contextError with a stack trace too.
func Errorf(format string, args ...interface{}) error {
	return &contextError{
		err:        nil,
		ctx:        fmt.Sprintf(format, args...),
		stacktrace: getStackTrace(),
	}
}

// contextError implements the error interface. It contains an error, context
// string, and a stack trace.
type contextError struct {
	err error
	ctx string
	*stacktrace
}

// Error allows ContextError to implement the error interface.
func (e *contextError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%v: %s", e.ctx, e.err.Error())
	}
	return e.ctx
}

// As implements the As interface as specified in the errors library.
// See https://golang.org/pkg/errors/#As
// TODO(terrenceho): Figure out what functionality we actually want from As().
func (e *contextError) As(target interface{}) bool {
	_, ok := target.(*contextError)
	if !ok {
		return false
	}
	target = e
	return true
}

// Is implements the Is interface as specified in the errors library.
// See https://golang.org/pkg/errors/#Is
func (e *contextError) Is(target error) bool {
	c_target, ok := target.(*contextError)
	if ok && e.err != nil && c_target.err != nil {
		return e.ctx == c_target.ctx && e.err.Error() == c_target.err.Error()
	}
	return target.Error() == e.Error()
}

// Unwrap implements stderror's Unwrap function. Returns the previous error
// without context.
// See https://golang.org/pkg/errors/#Unwrap
func (e *contextError) Unwrap() error {
	return e.err
}

// Wrap is a convenience function around Wrapf
func Wrap(err error, context string) error {
	return &contextError{
		err:        err,
		ctx:        context,
		stacktrace: getStackTrace(),
	}
}

// Wrapf takes in an error, a string with optional formatting options, and
// returns an error that wraps the input error. The returned error can be
// unwrapped using the errors.Unwrap() method.
//
// Under the hood, Wrapf returns a ContextError that stores both the original
// error and the formatted string.
func Wrapf(err error, format string, args ...interface{}) error {
	return &contextError{
		err:        err,
		ctx:        fmt.Sprintf(format, args...),
		stacktrace: getStackTrace(),
	}
}

// RootCause returns the root error that caused the chain or errors to being
// with. It repeatedly tries to unwrap the error until it can no longer be
// unwrapped.
func RootCause(err error) error {
	for {
		ctx_err, ok := err.(*contextError)
		if !ok {
			return err
		}
		if ctx_err.err == nil {
			return New(ctx_err.ctx)
		}
		err = ctx_err.err
	}
}

///// formatting for contextError

// Format implements the Formatter interface. It supports the following 3 verbs.
//
//		%s	print the error. If it is a contextError, the Error() implementation
//			will cause the error to print recursively down the chain.
//		%v	same as %s.
//		%q  prints the error but in quotes, according to go's verb usage.
//
// The %v also has another option: %+v. This will print the stacktrace in the
// following format
//
//		error message 3
//		--- /path/github.com/repo/package/main.go:<line num> (FuncName)
//		error message 2
//		--- /path/github.com/repo/package/main.go:<line num> (Class.FuncName)
//		error message 1
//		--- /path/github.com/repo/package/main.go:<line num> ((*Class).FuncName)
//
// All other formating cases will fall through.
func (ctx *contextError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			msg := ctx.fullStackTrace("", s, verb)
			fmt.Fprintf(s, "%s", msg)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, ctx.Error())
	case 'q':
		fmt.Fprintf(s, "%q", ctx.Error())
	}
}

// recursively builds stack trace and returns the full formatted string.
func (ctx *contextError) fullStackTrace(msg string, s fmt.State, verb rune) string {
	msg = addNewline(msg)
	msg += fmt.Sprint(ctx.ctx)
	msg = addNewline(msg)
	msg += fmt.Sprintf("%+v", ctx.stacktrace)
	msg = addNewline(msg)

	// recursively call format to print stack trace.
	if ctx.err != nil {
		f_ctx, ok := (ctx.err).(*contextError)
		if ok {
			msg += fmt.Sprintf("%+v", f_ctx)
		}
	}
	return msg
}

// addNewLine adds a new line if it doesn't end in a newline or is not empty
func addNewline(msg string) string {
	if msg != "" && !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return msg
}
