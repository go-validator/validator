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

package walidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// TextErr is an error that also implements the encoding.TextMarshaler interface for
// serializing out to various plain text encodings. Packages creating their
// own custom errors should use TextErr if they're intending to use serializing
// formats like json, msgpack etc.
type TextErr struct {
	Err error
}

// Error implements the error interface.
func (t TextErr) Error() string {
	return t.Err.Error()
}

// MarshalText implements encoding.TextMarshaler
func (t TextErr) MarshalText() ([]byte, error) {
	return []byte(t.Err.Error()), nil
}

var (
	// ErrZeroValue is the error returned when variable has zero valud
	// and nonzero was specified
	ErrZeroValue = TextErr{errors.New("zero value")}
	// ErrMin is the error returned when variable is less than mininum
	// value specified
	ErrMin = TextErr{errors.New("less than min")}
	// ErrMax is the error returned when variable is more than
	// maximum specified
	ErrMax = TextErr{errors.New("greater than max")}
	// ErrLen is the error returned when length is not equal to
	// param specified
	ErrLen = TextErr{errors.New("invalid length")}
	// ErrRegexp is the error returned when the value does not
	// match the provided regular expression parameter
	ErrRegexp = TextErr{errors.New("regular expression mismatch")}
	// ErrUnsupported is the error error returned when a validation rule
	// is used with an unsupported variable type
	ErrUnsupported = TextErr{errors.New("unsupported type")}
	// ErrBadParameter is the error returned when an invalid parameter
	// is provided to a validation rule (e.g. a string where an int was
	// expected (max=foo,len=bar) or missing a parameter when one is required (len=))
	ErrBadParameter = TextErr{errors.New("bad parameter")}
	// ErrUnknownTag is the error returned when an unknown tag is found
	ErrUnknownTag = TextErr{errors.New("unknown tag")}
	// ErrInvalid is the error returned when variable is invalid
	// (normally a nil pointer)
	ErrInvalid = TextErr{errors.New("invalid value")}
	// ErrRequired is the error returned when variable is nil and
	// required tag was specified
	ErrRequired = TextErr{errors.New("required value")}
)

// ErrorMap is a map which contains all errors from validating a struct.
type ErrorMap map[string]ErrorArray

// ErrorMap implements the Error interface so we can check error against nil.
// The returned error is if existent the first error which was added to the map.
func (err ErrorMap) Error() string {
	for k, errs := range err {
		if len(errs) > 0 {
			return fmt.Sprintf("%s: %s", k, errs.Error())
		}
	}
	return ""
}

// ErrorArray is a slice of errors returned by the Validate function.
type ErrorArray []error

// ErrorArray implements the Error interface and returns the first error as
// string if existent.
func (err ErrorArray) Error() string {
	if len(err) > 0 {
		return err[0].Error()
	}
	return ""
}

// validationFunc is the internal form of a validation function.
// It validates the given value, adding any validation errors
// to the validation state.
type validationFunc func(reflect.Value, *validateState)

// okValidation is the no-op validator - it always succeeds.
func okValidation(reflect.Value, *validateState) {}

// errorValidation always fails with the given error.
func errorValidation(err error) validationFunc {
	return func(_ reflect.Value, state *validateState) {
		state.error(err)
	}
}

// ValidationFunc is a function that receives the value of a
// field and a parameter used for the respective validation tag.
type ValidationFunc func(v interface{}, param string) error

// Validator implements a validator
type Validator struct {
	// Tag name being used.
	tagName string
	// validationFuncs is a map of ValidationFuncs indexed
	// by their name.
	validationFuncs map[string]tagValidator

	validatorCache sync.Map // map[reflect.Type]validationFunc
}

// Helper validator so users can use the
// functions directly from the package
var defaultValidator = NewValidator()

// NewValidator creates a new Validator
func NewValidator() *Validator {
	return &Validator{
		tagName: "validate",
		validationFuncs: map[string]tagValidator{
			"nonzero":   legacyTagValidator(nonzero),
			"len":       legacyTagValidator(length),
			"min":       legacyTagValidator(min),
			"max":       legacyTagValidator(max),
			"regexp":    legacyTagValidator(regex),
			"uuid":      legacyTagValidator(uuid),
			"required":  legacyTagValidator(required),
			"latitude":  legacyTagValidator(latitude),
			"longitude": legacyTagValidator(longitude),
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
	// Setting the tag invalidates the cache.
	mv.validatorCache = sync.Map{}
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
	mv1 := &Validator{
		tagName:         mv.tagName,
		validationFuncs: make(map[string]tagValidator),
	}
	for k, f := range mv.validationFuncs {
		mv1.validationFuncs[k] = f
	}
	return mv1
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
	mv.validationFuncs[name] = legacyTagValidator(vf)
	return nil
}

// Validate validates the fields of a struct based
// on 'validate' tags and returns errors found indexed
// by the field name.
func Validate(v interface{}) error {
	return defaultValidator.Validate(v)
}

// Validate validates the fields of a struct based
// on 'validator' tags and returns errors found indexed
// by the field name.
func (mv *Validator) Validate(x interface{}) error {
	sv := reflect.ValueOf(x)
	validate := mv.typeValidator(sv.Type())
	// TODO calculate likely size of path and pathStack; or alternatively
	// we could keep validateState instances around in a sync.Pool
	// to avoid the allocations.
	state := &validateState{
		path:      make([]byte, 0, 20),
		pathStack: make([]int, 0, 10),
	}
	validate(sv, state)
	return state.finalError()
}

// jsonFieldName returns the name that the field will be given
// when marshaled to JSON, or the empty string if
func jsonFieldName(tag reflect.StructTag) string {
	jtag := tag.Get("json")
	if jtag == "" || jtag == "-" {
		return ""
	}
	i := strings.Index(jtag, ",")
	if i >= 0 {
		return jtag[0:i]
	}
	return jtag
}

// Valid validates a value based on the provided
// tags and returns errors found or nil.
func Valid(val interface{}, tags string) error {
	return defaultValidator.Valid(val, tags)
}

// Valid validates a value based on the provided
// tags and returns errors found or nil.
func (mv *Validator) Valid(v interface{}, tag string) error {
	if tag == "-" {
		return nil
	}
	sv := reflect.ValueOf(v)
	var st reflect.Type
	if sv.IsValid() {
		st = sv.Type()
	}
	validate, err := mv.parseTags(tag, st)
	if err != nil {
		// unknown tag found, give up.
		return err
	}
	state := &validateState{
		path:      make([]byte, 0, 20),
		pathStack: make([]int, 0, 10),
	}
	validate(sv, state)
	if err := state.finalError(); err != nil {
		// For backward compatibility, we're expected to return an ErrorArray here.
		if _, ok := err.(ErrorArray); ok {
			return err
		}
		return ErrorArray{err}
	}
	return nil
}

// newTypeValidator returns a validation function that checks values of type t.
func (mv *Validator) newTypeValidator(t reflect.Type) validationFunc {
	switch t.Kind() {
	case reflect.Ptr:
		elemf := mv.typeValidator(t.Elem())
		return func(v reflect.Value, state *validateState) {
			if v.IsNil() {
				return
			}
			elemf(v.Elem(), state)
		}
	case reflect.Struct:
		return mv.newStructValidator(t)
	case reflect.Array, reflect.Slice:
		elemf := mv.typeValidator(t.Elem())
		return func(v reflect.Value, state *validateState) {
			n := v.Len()
			for i := 0; i < n; i++ {
				state.pushPathIndex(i)
				elemf(v.Index(i), state)
				state.popPath()
			}
		}
	case reflect.Interface:
		return func(v reflect.Value, state *validateState) {
			if v.IsNil() {
				return
			}
			iv := v.Elem()
			mv.typeValidator(iv.Type())(iv, state)
		}
	case reflect.Map:
		keyf := mv.typeValidator(t.Key())
		elemf := mv.typeValidator(t.Elem())
		return func(v reflect.Value, state *validateState) {
			iter := v.MapRange()
			for iter.Next() {
				mk := iter.Key()
				state.pushPathMapKey(mk)
				keyf(mk, state)
				state.popPath()
				mv := iter.Value()
				state.pushPathMapVal(mk)
				elemf(mv, state)
				state.popPath()
			}
		}
	default:
		return okValidation
	}
}

type field struct {
	index    []int
	name     string
	validate validationFunc
}

type structValidator struct {
	fields []field
}

func (s *structValidator) validate(v reflect.Value, ectx *validateState) {
	for i := range s.fields {
		f := &s.fields[i]
		ectx.pushPathField(f.name)
		f.validate(v.FieldByIndex(f.index), ectx)
		ectx.popPath()
	}
}

func (mv *Validator) newStructValidator(t reflect.Type) validationFunc {
	var sv structValidator
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		name := f.Name
		tag := f.Tag.Get(mv.tagName)
		if tag == "-" {
			continue
		}
		if jsonName := jsonFieldName(f.Tag); jsonName != "" {
			name = jsonName
		}
		tagValidator, err := mv.parseTags(tag, f.Type)
		if err != nil {
			sv.fields = append(sv.fields, field{
				index:    f.Index,
				name:     name,
				validate: errorValidation(err),
			})
			continue
		}
		fieldValidator := mv.typeValidator(f.Type)
		sv.fields = append(sv.fields, field{
			index: f.Index,
			name:  name,
			validate: func(v reflect.Value, state *validateState) {
				tagValidator(v, state)
				fieldValidator(v, state)
			},
		})
	}
	return sv.validate
}

// typeValidator is like newTypeValidator except that it returns
// a cached validation function if possible.
func (mv *Validator) typeValidator(t reflect.Type) validationFunc {
	if vf, ok := mv.validatorCache.Load(t); ok {
		return vf.(validationFunc)
	}
	// Stolen logic from encoding/json...
	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  validationFunc
	)
	wg.Add(1)
	fi, loaded := mv.validatorCache.LoadOrStore(t, validationFunc(func(v reflect.Value, state *validateState) {
		wg.Wait()
		f(v, state)
	}))
	if loaded {
		return fi.(validationFunc)
	}

	// Compute the real encoder and replace the indirect func with it.
	f = mv.newTypeValidator(t)
	wg.Done()
	mv.validatorCache.Store(t, f)
	return f
}

// separate by no escaped commas
var sepPattern = regexp.MustCompile(`((?:^|[^\\])(?:\\\\)*),`)

func splitUnescapedComma(str string) []string {
	var ret []string
	indexes := sepPattern.FindAllStringIndex(str, -1)
	last := 0
	for _, is := range indexes {
		ret = append(ret, str[last:is[1]-1])
		last = is[1]
	}
	ret = append(ret, str[last:])
	return ret
}

// parseTags returns a validation function that checks the validations
// implied by the struct tag t that tags a field with the given type.
func (mv *Validator) parseTags(t string, fieldType reflect.Type) (validationFunc, error) {
	if t == "" {
		return okValidation, nil
	}
	tl := splitUnescapedComma(t)
	var validators []validationFunc
	for _, i := range tl {
		i = strings.Replace(i, `\,`, ",", -1)
		v := strings.SplitN(i, "=", 2)
		name := strings.Trim(v[0], " ")
		if name == "" {
			return nil, ErrUnknownTag
		}
		var param string
		if len(v) > 1 {
			param = strings.Trim(v[1], " ")
		}
		tvf, ok := mv.validationFuncs[name]
		if !ok {
			return nil, ErrUnknownTag
		}
		vf, err := tvf(fieldType, param)
		if err != nil {
			return nil, err
		}
		validators = append(validators, vf)
	}
	if len(validators) == 0 {
		return okValidation, nil
	}
	return func(v reflect.Value, state *validateState) {
		for _, f := range validators {
			f(v, state)
		}
	}, nil
}

type tagValidator func(t reflect.Type, param string) (validationFunc, error)

// legacyTagValidator converts from an external ValidationFunc
// to the internal form.
func legacyTagValidator(f ValidationFunc) tagValidator {
	return func(t reflect.Type, param string) (validationFunc, error) {
		return func(v reflect.Value, state *validateState) {
			var iv interface{}
			if v.IsValid() {
				iv = v.Interface()
			}
			if err := f(iv, param); err != nil {
				state.error(err)
			}
		}, nil
	}
}

// validateState holds the runtime state maintained when validating
// a value. It uses a single byte slice for the path to avoid allocations
// in the frequently used happy path when no errors are created.
//
// All calls to the push* methods must be balanced by popPath calls.
type validateState struct {
	path      []byte
	pathStack []int
	errors    ErrorMap
}

// finalError returns an error value that includes all the errors
// added by state.error.
func (state *validateState) finalError() error {
	if state.errors == nil {
		return nil
	}
	if len(state.errors) == 1 && len(state.errors[""]) > 0 {
		errs := state.errors[""]
		if len(errs) == 1 {
			return errs[0]
		}
		return errs
	}
	return state.errors
}

// error adds the given error to the set of errors recorded in state.
func (state *validateState) error(err error) {
	name := string(state.path)
	if state.errors == nil {
		state.errors = make(ErrorMap)
	}
	state.errors[name] = append(state.errors[name], err)
}

// pushPathField pushes a field name onto the current path.
func (state *validateState) pushPathField(fieldName string) {
	state._pushPath()
	if len(state.path) > 0 {
		state.path = append(state.path, '.')
	}
	state.path = append(state.path, []byte(fieldName)...)
}

// pushPathField pushes a map key onto the current path.
func (state *validateState) pushPathMapKey(key reflect.Value) {
	state._pushPath()
	state.path = append(state.path, []byte(fmt.Sprintf("[%+v](key)", key))...)
}

// pushPathField pushes a map value onto the current path.
func (state *validateState) pushPathMapVal(key reflect.Value) {
	state._pushPath()
	state.path = append(state.path, []byte(fmt.Sprintf("[%+v](value)", key))...)
}

// pushPathField pushes a slice or array index onto the current path.
func (state *validateState) pushPathIndex(i int) {
	state._pushPath()
	state.path = append(state.path, '[')
	state.path = strconv.AppendInt(state.path, int64(i), 10)
	state.path = append(state.path, ']')
}

// popPath undoes the most recent push* call.
func (state *validateState) popPath() {
	i := len(state.pathStack) - 1
	state.path = state.path[0:state.pathStack[i]]
	state.pathStack = state.pathStack[0:i]
}

func (state *validateState) _pushPath() {
	state.pathStack = append(state.pathStack, len(state.path))
}
