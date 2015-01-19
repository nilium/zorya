package zorya

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

const traceSize = 64

const (
	CodeNoErr int = iota
	CodeUnderflow
	CodeOverflow
	CodeBadArg
	CodeBadAccess
	CodeBadOpcode
)

type Error struct {
	Code    int
	Msg     string
	callers []uintptr
}

func (e *Error) Error() string {
	return e.Msg
}

// reify allocates a new Error with the same Code and Msg, but with a new trace
// slice for the current Callers above the reify call.
func (e *Error) reify() *Error {
	var buf [traceSize]uintptr
	var bufsl = buf[:]
	numCallers := runtime.Callers(2, bufsl)
	var exact = make([]uintptr, numCallers)
	copy(exact, bufsl)
	return &Error{e.Code, e.Msg, exact}
}

// Trace returns an array of strings where each element in the array describes
// the recorded for the error.
func (e *Error) Trace() []string {
	calls := make([]string, len(e.callers))
	for i, pc := range e.callers {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			calls[i] = "<unknown>"
			continue
		}
		name := fn.Name()
		if name != "runtime.sigpanic" {
			pc -= 1
		}
		file, line := fn.FileLine(pc)
		entry := fn.Entry()

		file = filepath.Base(file)
		calls[i] = fmt.Sprintf("[%5d:%012x] %v => %v (%012x)", line, pc, file, name, entry)
	}
	return calls
}

// TraceStr returns the stack trace for the Error as a single string rather
// than an array of strings.
func (e *Error) TraceStr() string {
	return strings.Join(e.Trace(), "\n")
}

// IsBadOpcode is used to determine whether a particular error represents a
// bad-opcode error in Zorya.
func IsBadOpcode(err error) bool {
	return HasCode(err, CodeBadOpcode)
}

// IsBadAccess is used to determine whether a particular error represents a
// bad-access error in Zorya.
func IsBadAccess(err error) bool {
	return HasCode(err, CodeBadAccess)
}

// ErrorCode returns the error code for a Zorya Error if err is a Zorya error,
// otherwise returns CodeNoErr.
func ErrorCode(err error) int {
	e, ok := err.(*Error)
	if !ok {
		return CodeNoErr
	}
	return e.Code
}

// HasCode returns whether the given error, if it is an Error, has a specific
// error code. If err is not an Error, this will return true only if code is
// CodeNoErr (which does not indicate that the error is nil, only that Zorya
// does not recognize it).
//
// This can be used to identify and respond to specific classes of error in
// Zorya, though in practice, if an error has occurred in Zorya, a thread's
// state may be undefined.
func HasCode(err error, code int) bool {
	return ErrorCode(err) == code
}

var (
	errStackUnderflow = &Error{CodeUnderflow, "Stack underflow", nil}
	errStackOverflow  = &Error{CodeOverflow, "Stack overflow", nil}
)

// mkerrorf returns a new Error with the given code, format string, and applied
// format parameters. Its Trace field is automatically populated with Callers
// above the mkerrorf call.
func mkerrorf(code int, format string, args ...interface{}) *Error {
	var buf [traceSize]uintptr
	var bufsl = buf[:]
	numCallers := runtime.Callers(2, bufsl)
	var exact = make([]uintptr, numCallers)
	copy(exact, bufsl)

	return &Error{
		Code:    code,
		Msg:     fmt.Sprintf(format, args...),
		callers: exact,
	}
}
