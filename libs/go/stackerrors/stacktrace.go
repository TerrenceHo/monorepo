package stackerrors

import (
	"fmt"
	"runtime"
	"strings"
)

type stacktrace struct {
	file     string // File where the error originated
	function string // function or method where the error originated
	line     int    // line number where the error originated
}

func (st *stacktrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			msg := formatStackEntry(st)
			fmt.Fprintf(s, msg)
			return
		}
	}
}

func formatStackEntry(st *stacktrace) string {
	var msg string
	msg = fmt.Sprintf("--- %s:%d (%s) ---", st.file, st.line, st.function)
	return msg
}

// getStackTrace returns a stacktrace of the function that created the error.
func getStackTrace() *stacktrace {
	// getStackTrace is called from Wrapf(), New(), or Errorf(), so two levels
	// above the caller will be the function that created the error
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		// Debug symbols may have been stripped from the binary, so
		// runtime.Caller may fail to return anything.
		return nil
	}
	funcName := funcFromPC(pc)

	st := stacktrace{
		file:     file,
		line:     line,
		function: funcName,
	}

	return &st
}

// funcFromPC takes in a pc counter and returns the function name without a path
// as a string. There are three possible types of function paths:
// - github.com/org/path/package.FuncName (Function)
// - github.com/org/path/package.Class.FuncName (Method)
// - github.com/org/path/package.(*Class).FuncName (Pointer Method)
func funcFromPC(pc uintptr) string {
	f := runtime.FuncForPC(pc)
	funcName := f.Name()
	return processFuncName(funcName)
}

func processFuncName(funcName string) string {
	// remove path parts by indexing to the last "/"
	funcName = funcName[strings.LastIndex(funcName, "/")+1:]

	// remove "package" name
	funcName = funcName[strings.Index(funcName, ".")+1:]

	return funcName
}
