package platform

import "fmt"

type Code string

const (
	CodeInvalid      Code = "invalid_request"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeUnauthorized Code = "unauthorized"
	CodeInternal     Code = "internal"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
func (e *Error) Unwrap() error             { return nil }
func Invalid(msg string, err error) error  { return &Error{Code: CodeInvalid, Message: msg, Err: err} }
func NotFound(msg string) error            { return &Error{Code: CodeNotFound, Message: msg} }
func Conflict(msg string) error            { return &Error{Code: CodeConflict, Message: msg} }
func Internal(msg string, err error) error { return &Error{Code: CodeInternal, Message: msg, Err: err} }
