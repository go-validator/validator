// Package validator implements value validations
//
// Copyright 2014 Roberto Teixeira <robteix@robteix.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator

import (
	"errors"
	"reflect"
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

// ValidationFunc is a function that receives the value of a
// field and a parameter used for the respective validation tag.
type ValidationFunc func(v interface{}, param string) error

// Validator implements a validator
type Validator struct {
	// Tag name being used.
	tagName string
	// validationFuncs is a map of ValidationFuncs indexed
	// by their name.
	validationFuncs map[string]ValidationFunc
}

// Helper validator so users can use the
// functions directly from the package
var defaultValidator = NewValidator()

// NewValidator creates a new Validator
func NewValidator() *Validator {
	return &Validator{
		tagName: "validate",
		validationFuncs: map[string]ValidationFunc{
			"nonzero": nonzero,
			"len":     length,
			"min":     min,
			"max":     max,
			"regexp":  regex,
		},
	}
}

// SetTag allows you to change the tag name used in structs
func SetTag(tag string) {
	defaultValidator.SetTag(tag)
}

// SetTag allows you to change the tag name used in structs
func (mv *Validator) SetTag(tag string) {
	mv.tagName = tag
}

// WithTag creates a new Validator with the new tag name. It is
// useful to chain-call with Validate so we don't change the tag
// name permanently: validator.WithTag("foo").Validate(t)
func WithTag(tag string) *Validator {
	return defaultValidator.WithTag(tag)
}

// WithTag creates a new Validator with the new tag name. It is
// useful to chain-call with Validate so we don't change the tag
// name permanently: validator.WithTag("foo").Validate(t)
func (mv *Validator) WithTag(tag string) *Validator {
	v := mv.copy()
	v.SetTag(tag)
	return v
}

// Copy a validator
func (mv *Validator) copy() *Validator {
	return &Validator{
		tagName:         mv.tagName,
		validationFuncs: mv.validationFuncs,
	}
}

// SetValidationFunc sets the function to be used for a given
// validation constraint. Calling this function with nil vf
// is the same as removing the constraint function from the list.
func SetValidationFunc(name string, vf ValidationFunc) error {
	return defaultValidator.SetValidationFunc(name, vf)
}

// SetValidationFunc sets the function to be used for a given
// validation constraint. Calling this function with nil vf
// is the same as removing the constraint function from the list.
func (mv *Validator) SetValidationFunc(name string, vf ValidationFunc) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if vf == nil {
		delete(mv.validationFuncs, name)
		return nil
	}
	mv.validationFuncs[name] = vf
	return nil
}

// Validate validates the fields of a struct based
// on 'validator' tags and returns errors found indexed
// by the field name.
func Validate(v interface{}) (bool, map[string][]error) {
	return defaultValidator.Validate(v)
}

// Validate validates the fields of a struct based
// on 'validator' tags and returns errors found indexed
// by the field name.
func (mv *Validator) Validate(v interface{}) (bool, map[string][]error) {
	sv := reflect.ValueOf(v)
	st := reflect.TypeOf(v)
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		return mv.Validate(sv.Elem().Interface())
	}
	if sv.Kind() != reflect.Struct {
		return false, map[string][]error{
			"": []error{ErrUnsupported},
		}
	}

	nfields := sv.NumField()
	m := make(map[string][]error)
	for i := 0; i < nfields; i++ {
		f := sv.Field(i)
		// deal with pointers
		for f.Kind() == reflect.Ptr && !f.IsNil() {
			f = f.Elem()
		}
		tag := st.Field(i).Tag.Get(mv.tagName)
		if f.Kind() == reflect.Ptr {
			ff := f.Elem()
			if ff.Kind() == reflect.Struct {

			}
		}
		if tag == "" && f.Kind() != reflect.Struct {
			continue
		}
		fname := st.Field(i).Name
		var errs []error
		switch f.Kind() {
		case reflect.Struct:
			_, e := mv.Validate(f.Interface())
			if len(e) > 0 {
				for j, k := range e {
					m[fname+"."+j] = k
				}
			}

		default:
			_, errs = mv.Valid(f.Interface(), tag)
		}

		if len(errs) > 0 {
			m[st.Field(i).Name] = errs
		}
	}
	return len(m) == 0, m
}

// Valid validates a value based on the provided
// tags and returns errors found or nil.
func Valid(val interface{}, tags string) (bool, []error) {
	return defaultValidator.Valid(val, tags)
}

// Valid validates a value based on the provided
// tags and returns errors found or nil.
func (mv *Validator) Valid(val interface{}, tags string) (bool, []error) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		return mv.Valid(v.Elem().Interface(), tags)
	}
	var errs []error
	switch v.Kind() {
	case reflect.Struct:
		return false, []error{ErrUnsupported}
	case reflect.Invalid:
		errs = mv.validateVar(nil, tags)
	default:
		errs = mv.validateVar(val, tags)
	}
	return len(errs) < 1, errs
}

// validateVar validates one single variable
func (mv *Validator) validateVar(v interface{}, tag string) []error {
	tags, err := mv.parseTags(tag)
	if err != nil {
		// unknown tag found, give up.
		return []error{err}
	}
	errs := make([]error, 0, len(tags))
	for _, t := range tags {
		if err := t.Fn(v, t.Param); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// tag represents one of the tag items
type tag struct {
	Name  string         // name of the tag
	Fn    ValidationFunc // validation function to call
	Param string         // parameter to send to the validation function
}

// parseTags parses all individual tags found within a struct tag.
func (mv *Validator) parseTags(t string) ([]tag, error) {
	tl := strings.Split(t, ",")
	tags := make([]tag, 0, len(tl))
	for _, i := range tl {
		tg := tag{}
		v := strings.SplitN(i, "=", 2)
		tg.Name = strings.Trim(v[0], " ")
		if tg.Name == "" {
			return []tag{}, ErrUnknownTag
		}
		if len(v) > 1 {
			tg.Param = strings.Trim(v[1], " ")
		}
		var found bool
		if tg.Fn, found = mv.validationFuncs[tg.Name]; !found {
			return []tag{}, ErrUnknownTag
		}
		tags = append(tags, tg)

	}
	return tags, nil
}
