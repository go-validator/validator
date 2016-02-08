// Package validator implements value validations
//
// Copyright (C) 2014-2016 Roberto Teixeira <robteix@robteix.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package validator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrZeroValue is the error returned when variable has zero valud
	// and nonzero was specified
	ErrZeroValue = errors.New("zero value")
	// ErrMin is the error returned when variable is less than mininum
	// value specified
	ErrMin = errors.New("less than min")
	// ErrMax is the error returned when variable is more than
	// maximum specified
	ErrMax = errors.New("greater than max")
	// ErrLen is the error returned when length is not equal to
	// param specified
	ErrLen = errors.New("invalid length")
	// ErrRegexp is the error returned when the value does not
	// match the provided regular expression parameter
	ErrRegexp = errors.New("regular expression mismatch")
	// ErrUnsupported is the error error returned when a validation rule
	// is used with an unsupported variable type
	ErrUnsupported = errors.New("unsupported type")
	// ErrBadParameter is the error returned when an invalid parameter
	// is provided to a validation rule (e.g. a string where an int was
	// expected (max=foo,len=bar) or missing a parameter when one is required (len=))
	ErrBadParameter = errors.New("bad parameter")
	// ErrUnknownTag is the error returned when an unknown tag is found
	ErrUnknownTag = errors.New("unknown tag")
	// ErrInvalid is the error returned when variable is invalid
	// (normally a nil pointer)
	ErrInvalid = errors.New("invalid value")
)

// Error implements a custom error that includes all
// validation errors for a given type
type Error struct {
	errs map[string][]error
}

// newError creates a new Error
func newError() *Error {
	return &Error{
		errs: map[string][]error{},
	}
}

// Errors returns all errors for a given field
func Errors(err error, field string) []error {
	verr, ok := err.(*Error)
	if !ok {
		return nil
	}
	return verr.errs[field]
}

// ErrorFields returns all fields with errors in err
func ErrorFields(err error) []string {
	verr, ok := err.(*Error)
	if !ok {
		return nil
	}
	fs := []string{}
	for i := range verr.errs {
		fs = append(fs, i)
	}
	return fs
}

// Error returns a string representation of the errors.
func (e Error) Error() string {
	errs := []string{}
	for k, v := range e.errs {
		for _, err := range v {
			errs = append(errs, fmt.Sprintf("%s has error %v", k, err))
		}
	}
	return strings.Join(errs, ", ")
}

// isError checks whether err for
// a given field is fieldError.
func isError(err error, field string, fieldError error) bool {
	verr, ok := err.(*Error)
	if !ok {
		return false
	}
	for _, e := range verr.errs[field] {
		if e == fieldError {
			return true
		}
	}
	return false
}

// IsZeroValue checks if err is ErrZeroValue for
// a given field.
func IsZeroValue(err error, field string) bool {
	return isError(err, field, ErrZeroValue)
}

// IsMin checks if err is ErrMin for
// a given field.
func IsMin(err error, field string) bool {
	return isError(err, field, ErrMin)
}

// IsMax checks if err is ErrMax for
// a given field.
func IsMax(err error, field string) bool {
	return isError(err, field, ErrMax)
}

// IsLen checks if err is ErrLen for
// a given field.
func IsLen(err error, field string) bool {
	return isError(err, field, ErrLen)
}

// IsRegexp checks if err is ErrRegexp for
// a given field.
func IsRegexp(err error, field string) bool {
	return isError(err, field, ErrRegexp)
}

// IsUnsupported checks if err is ErrUnsupported for
// a given field.
func IsUnsupported(err error, field string) bool {
	return isError(err, field, ErrUnsupported)
}

// IsBadParameter checks if err is ErrBadParameter for
// a given field.
func IsBadParameter(err error, field string) bool {
	return isError(err, field, ErrBadParameter)
}

// IsUnknownTag checks if err is ErrUnknownTag for
// a given field.
func IsUnknownTag(err error, field string) bool {
	return isError(err, field, ErrUnknownTag)
}

// IsInvalid checks if err is ErrInvalid for
// a given field.
func IsInvalid(err error, field string) bool {
	return isError(err, field, ErrInvalid)
}
