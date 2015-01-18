package zorya

import "fmt"

const (
	CodeUnderflow int = iota
	CodeOverflow
	CodeBadAccess
	CodeBadOpcode
)

type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	return e.Msg
}

// IsBadOpcode is used to determine whether a particular error represents a
// bad-opcode error in Zorya.
func IsBadOpcode(err error) bool {
	if e, ok := err.(*Error); ok && e.Code == CodeBadOpcode {
		return true
	}
	return false
}

// IsBadAccess is used to determine whether a particular error represents a
// bad-access error in Zorya.
func IsBadAccess(err error) bool {
	if e, ok := err.(*Error); ok && e.Code == CodeBadAccess {
		return true
	}
	return false
}

var (
	ErrStackUnderflow = &Error{CodeUnderflow, "Stack underflow"}
	ErrStackOverflow  = &Error{CodeOverflow, "Stack overflow"}
)

func mkerrorf(code int, format string, args ...interface{}) *Error {
	return &Error{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
	}
}
