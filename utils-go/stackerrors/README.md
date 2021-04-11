# Stackerrors

Implementation of errors with stacktraces where errors are generated. Compatible
with go 1.13 errors package.

## Usage

You can create new errors with:

```
err := stackerrors.New("error string")
```

There is an `Errorf()` equivalent as well:

```
err := stackerrors.Errorf("format %s", string)
```

If you have a previously existing error and would like to wrap that error with a
context string, you can use the `Wrap()` or `Wrapf()` functions.

```
err := stackerrors.Wrap(prev_err, "context string")
// OR
err := stackerrors.Wrapf(prev_err, "format context %s", "string")
```

The `Errorf()` and `Wrapf()` can take all formatting options defined in the
`fmt` package. If you want to retrieve the first error, you can `RootCause()`
function.

The `New()`, `Errorf()`, `Wrap()`, and `Wrapf()` functions all capture a
stacktrace when the error was created. You can print the stacktrace using the
`%+v` formatting option.

```
func FuncName() {
        fmt.Printf("%+v", stackerrors.New("base error"))
        fmt.Println()
        err := stackerrors.New("first error")
        err = stackerrors.Wrap(err, "second error")
        err = stackerrors.Wrap(err, "third error")
        fmt.Printf("%+v", err)
}
```

Output:

```
base error
--- /file/path/to/program.go:<line number> (FuncName) ---

third error
--- /file/path/to/program.go:<line number> (FuncName) ---
second error
--- /file/path/to/program.go:<line number> (FuncName) ---
first error
--- /file/path/to/program.go:<line number> (FuncName) ---
```
