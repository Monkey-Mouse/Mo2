package mo2errors

import "fmt"

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

// Init as name
func (e *Mo2Errors) Init(c int, s string) {
	e.ErrorCode = c
	e.ErrorTip = s
}

// InitCode as name
func (e *Mo2Errors) InitCode(c int) {
	e.ErrorCode = c
	e.ErrorTip = CodeText(c)
}

// IsError as name
func (e Mo2Errors) IsError() (error bool) {
	error = true
	if e.ErrorCode == Mo2NoError {
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
