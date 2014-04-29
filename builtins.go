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
	"regexp"
	"strconv"
)

// nonzero tests whether a variable value non-zero
// as defined by the golang spec.
func nonzero(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	valid := true
	switch st.Kind() {
	case reflect.String:
		valid = len(st.String()) != 0
	case reflect.Ptr, reflect.Interface:
		valid = !st.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		valid = st.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = st.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid = st.Uint() != 0
	case reflect.Float32, reflect.Float64:
		valid = st.Float() != 0
	case reflect.Bool:
		valid = st.Bool()
	case reflect.Invalid:
		valid = false // always invalid
	default:
		panic("Unsupported type " + st.String())
	}

	if !valid {
		return ErrZeroValue
	}
	return nil
}

// length tests whether a variable's length is equal to a given
// value. For strings it tests the number of characters whereas
// for maps and slices it tests the number of items.
func length(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	valid := true
	switch st.Kind() {
	case reflect.String:
		valid = int64(len(st.String())) == asInt(param)
	case reflect.Slice, reflect.Map, reflect.Array:
		valid = int64(st.Len()) == asInt(param)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = st.Int() == asInt(param)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid = st.Uint() == asUint(param)
	case reflect.Float32, reflect.Float64:
		valid = st.Float() == asFloat(param)
	default:
		panic("length is not a valid validation tag for type " + st.String())
	}
	if !valid {
		return ErrLen
	}
	return nil
}

// min tests whether a variable value is larger or equal to a given
// number. For number types, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func min(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	invalid := false
	switch st.Kind() {
	case reflect.String:
		invalid = int64(len(st.String())) < asInt(param)
	case reflect.Slice, reflect.Map, reflect.Array:
		invalid = int64(st.Len()) < asInt(param)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		invalid = st.Int() < asInt(param)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		invalid = st.Uint() < asUint(param)
	case reflect.Float32, reflect.Float64:
		invalid = st.Float() < asFloat(param)
	default:
		panic("min is not a valid validation tag for type " + st.String())
	}
	if invalid {
		return ErrMin
	}
	return nil
}

// max tests whether a variable value is lesser than a given
// value. For numbers, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func max(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	var invalid bool
	switch st.Kind() {
	case reflect.String:
		invalid = int64(len(st.String())) > asInt(param)
	case reflect.Slice, reflect.Map, reflect.Array:
		invalid = int64(st.Len()) > asInt(param)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		invalid = st.Int() > asInt(param)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		invalid = st.Uint() > asUint(param)
	case reflect.Float32, reflect.Float64:
		invalid = st.Float() > asFloat(param)
	default:
		panic("max is not a valid validation tag for type " + st.String())
	}
	if invalid {
		return ErrMax
	}
	return nil
}

// regex is the builtin validation function that checks
// whether the string variable matches a regular expression
func regex(v interface{}, param string) error {
	s, ok := v.(string)
	if !ok {
		panic("regexp requires a string")
	}

	re := regexp.MustCompile(param)

	if !re.MatchString(s) {
		return errors.New("regexp does not match")
	}
	return nil
}

// asInt retuns the parameter as a int64
// or panics if it can't convert
func asInt(param string) int64 {
	i, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		panic("Invalid param " + param + ", should be an integer")
	}
	return i
}

// asUint retuns the parameter as a uint64
// or panics if it can't convert
func asUint(param string) uint64 {
	i, err := strconv.ParseUint(param, 0, 64)
	if err != nil {
		panic("Invalid param " + param + ", should be an unsigned integer")
	}
	return i
}

// asFloat retuns the parameter as a float64
// or panics if it can't convert
func asFloat(param string) float64 {
	i, err := strconv.ParseFloat(param, 64)
	if err != nil {
		panic("Invalid param " + param + ", should be a float")
	}
	return i
}
