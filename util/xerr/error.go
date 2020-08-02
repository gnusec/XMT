package xerr

// Error is a struct that can be used to fast-path wrap errors and prevent loading the stdlib error
// functions. The 'Wrap' function will append the additional error value's 'Error()' string to the
// end, if no nil. This function supports the 'Unwrap' function in the 'errors' package.
type Error struct {
	e error
	s string
}

// Error returns the error message of this Error as a string value.
func (e Error) Error() string {
	return e.s
}

// Unwrap supports the 'errors.Unwrap' function. This will return the wrapped error, if not nil.
func (e Error) Unwrap() error {
	return e.e
}

// String returns the string value of this Error. Similar to the 'Error()' function.
func (e Error) String() string {
	return e.s
}

// Wrap creates a new error that wraps the specified error. If not nil, this function will append ": " + 'Error()'
// to the resulting string message.
func Wrap(s string, e error) error {
	if e != nil {
		return &Error{s: s + ": " + e.Error(), e: e}
	}
	return &Error{s: s}
}