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
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Errors is the error type returned from the Validate function and method.
// The map key holds the path to where the validation failed,
// and each error in the associated slice holds the description of
// an issue that occurred at that path.
type Errors map[string][]Error

// Error holds a single validation error.
type Error struct {
	// Kind holds the name of the validator that failed
	// (for example "len"). If there was an error that
	// didn't stem from a particular validator (for example,
	// an unknown validator name), Kind will be "invalid".
	Kind string `json:"kind"`
	// Msg holds details about the failure.
	Msg string `json:"msg"`
}

// Error implements error.Error by returning the returning the
// first error found in the lowest valued map key (by lexical ordering).
func (errs Errors) Error() string {
	leastKey := ""
	var err *Error
	for k, errSlice := range errs {
		if len(errSlice) > 0 && (err == nil || k < leastKey) {
			leastKey = k
			err = &errSlice[0]
		}
	}
	if err == nil {
		// Defensive: shouldn't happen in practice.
		return "no validation errors"
	}
	if leastKey == "" {
		return "validation error: " + err.Msg
	}
	return "validation error: " + leastKey + ": " + err.Msg
}

// ValidationFunc validates a value, adding any validation
// errors by calling r.Errorf.
type ValidationFunc func(v reflect.Value, r *ErrorReporter)

// okValidation is the no-op validator - it always succeeds.
func okValidation(reflect.Value, *ErrorReporter) {}

// errorValidation always fails with the given error.
func errorValidation(err error) ValidationFunc {
	return func(_ reflect.Value, r *ErrorReporter) {
		r.Errorf("%v", err)
	}
}

// Validator implements a validator.
type Validator struct {
	// Tag name being used.
	tagName string
	// tagValidators is a map of TagValidators indexed by their name.
	tagValidators map[string]TagValidator

	validatorCache sync.Map // map[reflect.Type]ValidationFunc
}

// defaultValidator holds the validator used by the Validate function.
var defaultValidator = New()

// New creates a new Validator
func New() *Validator {
	return &Validator{
		tagName: "validate",
		tagValidators: map[string]TagValidator{
			"nonzero":   nonzero,
			"len":       concreteTagValidator(elemTagValidator(length)),
			"min":       concreteTagValidator(elemTagValidator(min)),
			"max":       concreteTagValidator(elemTagValidator(max)),
			"regexp":    concreteTagValidator(elemTagValidator(regex)),
			"uuid":      concreteTagValidator(elemTagValidator(uuid)),
			"required":  required,
			"latitude":  concreteTagValidator(elemTagValidator(latitude)),
			"longitude": concreteTagValidator(elemTagValidator(longitude)),
		},
	}
}

// WithTag creates a new Validator that's a copy of mv except that it
// recognises tags with the given tag name (the default is "validate").
func (mv *Validator) WithTag(tag string) *Validator {
	mv1 := &Validator{
		tagName:       mv.tagName,
		tagValidators: make(map[string]TagValidator),
	}
	for k, f := range mv.tagValidators {
		mv1.tagValidators[k] = f
	}
	mv1.tagName = tag
	return mv1
}

// AddValidation sets the validation function for the given
// name. If the name is empty or already registered, it panics.
func (mv *Validator) AddValidation(name string, vf TagValidator) {
	if name == "" {
		panic("empty validation name")
	}
	if vf == nil {
		panic("nil validation function")
	}
	if mv.tagValidators[name] != nil {
		panic("validation function " + name + " already registered")
	}
	mv.tagValidators[name] = vf
}

// Validate validates the given value. Tags in struct fields that have
// "validate" tags are checked according to the rules described in the
// package documentation.
//
// If there are validation errors, the returned error is of type Errors.
func Validate(v interface{}) error {
	return defaultValidator.Validate(v)
}

// Validate validates the given value. Tags in struct fields
// that have "validate" tags are checked according to the
// rules described in the package documentation,
// including additional validations added with mv.AddValidation.
//
// If there are validation errors, the returned error is of type Errors.
func (mv *Validator) Validate(x interface{}) error {
	if x == nil {
		return nil
	}
	sv := reflect.ValueOf(x)
	validate := mv.typeValidator(sv.Type())
	// TODO calculate likely size of path and pathStack; or alternatively
	// we could keep ErrorReporter instances around in a sync.Pool
	// to avoid the allocations.
	state := getState()
	validate(sv, state)
	err := state.finalError()
	putState(state)
	return err
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

// newTypeValidator returns a validation function that checks values of type t.
func (mv *Validator) newTypeValidator(t reflect.Type) ValidationFunc {
	switch t.Kind() {
	case reflect.Ptr:
		elemf := mv.typeValidator(t.Elem())
		return func(v reflect.Value, state *ErrorReporter) {
			if v.IsNil() {
				return
			}
			elemf(v.Elem(), state)
		}
	case reflect.Struct:
		return mv.newStructValidator(t)
	case reflect.Array, reflect.Slice:
		elemf := mv.typeValidator(t.Elem())
		return func(v reflect.Value, state *ErrorReporter) {
			n := v.Len()
			for i := 0; i < n; i++ {
				state.pushPathIndex(i)
				elemf(v.Index(i), state)
				state.popPath()
			}
		}
	case reflect.Interface:
		return func(v reflect.Value, state *ErrorReporter) {
			if v.IsNil() {
				return
			}
			iv := v.Elem()
			mv.typeValidator(iv.Type())(iv, state)
		}
	case reflect.Map:
		keyf := mv.typeValidator(t.Key())
		elemf := mv.typeValidator(t.Elem())
		return func(v reflect.Value, state *ErrorReporter) {
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
	validate ValidationFunc
}

type structValidator struct {
	fields []field
}

func (s *structValidator) validate(v reflect.Value, ectx *ErrorReporter) {
	for i := range s.fields {
		f := &s.fields[i]
		ectx.pushPathField(f.name)
		f.validate(v.FieldByIndex(f.index), ectx)
		ectx.popPath()
	}
}

func (mv *Validator) newStructValidator(t reflect.Type) ValidationFunc {
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
			validate: func(v reflect.Value, state *ErrorReporter) {
				tagValidator(v, state)
				fieldValidator(v, state)
			},
		})
	}
	return sv.validate
}

// typeValidator is like newTypeValidator except that it returns
// a cached validation function if possible.
func (mv *Validator) typeValidator(t reflect.Type) ValidationFunc {
	if vf, ok := mv.validatorCache.Load(t); ok {
		return vf.(ValidationFunc)
	}
	// Stolen logic from encoding/json...
	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  ValidationFunc
	)
	wg.Add(1)
	fi, loaded := mv.validatorCache.LoadOrStore(t, ValidationFunc(func(v reflect.Value, state *ErrorReporter) {
		wg.Wait()
		f(v, state)
	}))
	if loaded {
		return fi.(ValidationFunc)
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
func (mv *Validator) parseTags(t string, fieldType reflect.Type) (ValidationFunc, error) {
	if t == "" {
		return okValidation, nil
	}
	tl := splitUnescapedComma(t)
	var validators []ValidationFunc
	var kinds []string
	for _, i := range tl {
		i = strings.Replace(i, `\,`, ",", -1)
		v := strings.SplitN(i, "=", 2)
		name := strings.Trim(v[0], " ")
		if name == "" {
			return nil, fmt.Errorf("empty validation name in tag")
		}
		var param string
		if len(v) > 1 {
			param = strings.Trim(v[1], " ")
		}
		tvf, ok := mv.tagValidators[name]
		if !ok {
			return nil, fmt.Errorf("unknown validation name %q in tag", name)
		}
		vf, err := tvf(fieldType, param)
		if err != nil {
			vf = errorValidation(err)
		}
		validators = append(validators, vf)
		kinds = append(kinds, name)
	}
	if len(validators) == 0 {
		return okValidation, nil
	}
	return func(v reflect.Value, state *ErrorReporter) {
		for i, f := range validators {
			state.kind = kinds[i]
			f(v, state)
		}
		state.kind = ""
	}, nil
}

// TagValidator returns a validation function that will validate any values
// that have the given type. The param string holds the associated
// tag argument. The returned function will always be called with
// an argument of type t.
type TagValidator func(t reflect.Type, param string) (ValidationFunc, error)

// concreteTagValidator returns a tag validator that's the same
// as inner except that it ensures that inner will only see concrete
// (non-interface) values when it's called, by either unwrapping
// the interface value or by not calling inner for nil values.
func concreteTagValidator(inner TagValidator) TagValidator {
	// cache caches the resolved validation functions for
	// interface types.
	var cache sync.Map
	return func(t reflect.Type, param string) (ValidationFunc, error) {
		if t.Kind() != reflect.Interface {
			return inner(t, param)
		}
		return func(v reflect.Value, state *ErrorReporter) {
			v = v.Elem()
			if !v.IsValid() {
				return
			}
			t := v.Type()
			if innerf, ok := cache.Load(t); ok {
				innerf.(ValidationFunc)(v, state)
				return
			}
			innerf, err := inner(t, param)
			if err != nil {
				innerf = errorValidation(err)
			}
			cache.LoadOrStore(t, innerf)
			innerf(v, state)
		}, nil
	}
}

// elemTagValidator returns a tag validator that's the same as inner
// except that it will strip off one level of pointer indirection
// for pointer types.
func elemTagValidator(inner TagValidator) TagValidator {
	return func(t reflect.Type, param string) (ValidationFunc, error) {
		if t.Kind() != reflect.Ptr {
			return inner(t, param)
		}
		innerf, err := inner(t.Elem(), param)
		if err != nil {
			return nil, err
		}
		return func(v reflect.Value, state *ErrorReporter) {
			if v.IsNil() {
				return
			}
			innerf(v.Elem(), state)
		}, nil
	}
}

// ErrorReporter is used by custom validation functions to
// report validation errors.
type ErrorReporter struct {
	// kind holds the kind of validator currently being used.
	kind      string
	path      []byte
	pathStack []int
	errors    Errors
}

var reporterPool sync.Pool

func getState() *ErrorReporter {
	r, ok := reporterPool.Get().(*ErrorReporter)
	if ok {
		return r
	}
	return &ErrorReporter{
		path:      make([]byte, 0, 20),
		pathStack: make([]int, 0, 10),
	}
}

func putState(r *ErrorReporter) {
	r.path = r.path[:0]
	r.pathStack = r.pathStack[:0]
	r.errors = nil
	reporterPool.Put(r)
}

// finalError returns an error value that includes all the errors
// added by r.error.
func (r *ErrorReporter) finalError() error {
	if r.errors == nil {
		return nil
	}
	return r.errors
}

// Errorf adds the given error, formatted by fmt.Sprintf,
// to the set of errors recorded in r.
func (r *ErrorReporter) Errorf(format string, args ...interface{}) {
	name := string(r.path)
	if r.errors == nil {
		r.errors = make(Errors)
	}
	kind := r.kind
	if kind == "" {
		kind = "invalid"
	}
	r.errors[name] = append(r.errors[name], Error{
		Kind: kind,
		Msg:  fmt.Sprintf(format, args...),
	})
}

// pushPathField pushes a field name onto the current path.
func (r *ErrorReporter) pushPathField(fieldName string) {
	r._pushPath()
	if len(r.path) > 0 {
		r.path = append(r.path, '.')
	}
	r.path = append(r.path, []byte(fieldName)...)
}

// pushPathField pushes a map key onto the current path.
func (r *ErrorReporter) pushPathMapKey(key reflect.Value) {
	r._pushPath()
	r.path = append(r.path, []byte(fmt.Sprintf("[%+v](key)", key))...)
}

// pushPathField pushes a map value onto the current path.
func (r *ErrorReporter) pushPathMapVal(key reflect.Value) {
	r._pushPath()
	r.path = append(r.path, []byte(fmt.Sprintf("[%+v](value)", key))...)
}

// pushPathField pushes a slice or array index onto the current path.
func (r *ErrorReporter) pushPathIndex(i int) {
	r._pushPath()
	r.path = append(r.path, '[')
	r.path = strconv.AppendInt(r.path, int64(i), 10)
	r.path = append(r.path, ']')
}

// popPath undoes the most recent push* call.
func (r *ErrorReporter) popPath() {
	i := len(r.pathStack) - 1
	r.path = r.path[0:r.pathStack[i]]
	r.pathStack = r.pathStack[0:i]
}

func (r *ErrorReporter) _pushPath() {
	r.pathStack = append(r.pathStack, len(r.path))
}
