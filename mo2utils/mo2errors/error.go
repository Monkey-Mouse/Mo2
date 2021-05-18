package mo2errors

import (
	"fmt"
)

// Mo2Errors standard mo2 err
type Mo2Errors struct {
	ErrorCode int
	ErrorTip  string
}

func (e Mo2Errors) Error() string {
	return fmt.Sprintf("%v: %v", e.ErrorCode, e.ErrorTip)
}

// SetErrorTip as name
func (e *Mo2Errors) SetErrorTip(s string) {
	e.ErrorTip = s
}

// Init init with code and tip
func (e *Mo2Errors) Init(c int, format string, a ...interface{}) {
	e.ErrorCode = c
	e.ErrorTip = fmt.Sprintf(format, a...)
}

// Init init with code and tip
func Init(c int, format string, a ...interface{}) Mo2Errors {
	return Mo2Errors{
		ErrorCode: c,
		ErrorTip:  fmt.Sprintf(format, a...),
	}
}

// InitError init with error
func (e *Mo2Errors) InitError(err error) {
	e.ErrorCode = Mo2Error
	e.ErrorTip = err.Error()
}

// InitError init with error
func InitError(err error) Mo2Errors {
	return Mo2Errors{
		ErrorCode: Mo2Error,
		ErrorTip:  err.Error(),
	}
}

// InitNoError init with no error tips
func (e *Mo2Errors) InitNoError(format string, a ...interface{}) {
	e.ErrorCode = Mo2NoError
	e.ErrorTip = fmt.Sprintf(format, a...)
}

// InitNoError init with no error tips
func InitNoError(format string, a ...interface{}) Mo2Errors {
	return Mo2Errors{
		ErrorCode: Mo2NoError,
		ErrorTip:  fmt.Sprintf(format, a...),
	}
}

// InitCode init with code
func (e *Mo2Errors) InitCode(c int) {
	e.ErrorCode = c
	e.ErrorTip = CodeText(c)
}

// IsError as name
func (e Mo2Errors) IsError() (error bool) {
	error = true
	if e.ErrorCode == Mo2NoError || e == (Mo2Errors{}) {
		error = false
	}
	return
}

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(c int, s string) *Mo2Errors {
	return &Mo2Errors{c, s}
}

// NewCode returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func NewCode(c int) *Mo2Errors {
	return &Mo2Errors{c, CodeText(c)}
}
