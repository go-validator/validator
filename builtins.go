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
	"reflect"
	"regexp"
	"strconv"
)

// nonzero tests whether a variable value non-zero
// as defined by the golang spec.
func nonzero(t reflect.Type, param string) (validationFunc, error) {
	check := func(ok bool, state *validateState) {
		if !ok {
			state.error(ErrZeroValue)
		}
	}
	switch t.Kind() {
	case reflect.String:
		return func(v reflect.Value, state *validateState) {
			check(len(v.String()) != 0, state)
		}, nil
	case reflect.Ptr, reflect.Interface:
		return func(v reflect.Value, state *validateState) {
			check(!v.IsNil(), state)
		}, nil
	case reflect.Slice, reflect.Map, reflect.Array:
		return func(v reflect.Value, state *validateState) {
			check(v.Len() != 0, state)
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(v reflect.Value, state *validateState) {
			check(v.Int() != 0, state)
		}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return func(v reflect.Value, state *validateState) {
			check(v.Uint() != 0, state)
		}, nil
	case reflect.Float32, reflect.Float64:
		return func(v reflect.Value, state *validateState) {
			// TODO this preserves old behavior but is arguably
			// incorrect - it should probably error if the value is NaN.
			check(v.Float() != 0, state)
		}, nil
	case reflect.Bool:
		return func(v reflect.Value, state *validateState) {
			check(v.Bool(), state)
		}, nil
	case reflect.Struct:
		return okValidation, nil
	}
	return nil, ErrUnsupported
}

// length tests whether a variable's length is equal to a given
// value. For strings it tests the number of characters whereas
// for maps and slices it tests the number of items.
func length(t reflect.Type, param string) (validationFunc, error) {
	switch t.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if int64(v.Len()) != p {
				state.error(ErrLen)
			}
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if v.Int() != p {
				state.error(ErrLen)
			}
		}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := strconv.ParseUint(param, 0, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if v.Uint() != p {
				state.error(ErrLen)
			}
		}, nil
	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if v.Float() != p {
				state.error(ErrLen)
			}
		}, nil
	}
	return nil, ErrUnsupported
}

// min tests whether a variable value is larger or equal to a given
// number. For number types, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func min(t reflect.Type, param string) (validationFunc, error) {
	switch t.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if int64(v.Len()) < p {
				state.error(ErrMin)
			}
		}, nil
	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if v.Float() < p {
				state.error(ErrMin)
			}
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if v.Int() < p {
				state.error(ErrMin)
			}
		}, nil
	default:
		return nil, ErrUnsupported
	}
}

// max tests whether a variable value is lesser than a given
// value. For numbers, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func max(t reflect.Type, param string) (validationFunc, error) {
	switch t.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if int64(v.Len()) > p {
				state.error(ErrMax)
			}
		}, nil
	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if v.Float() > p {
				state.error(ErrMax)
			}
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, ErrBadParameter
		}
		return func(v reflect.Value, state *validateState) {
			if v.Int() > p {
				state.error(ErrMax)
			}
		}, nil
	default:
		return nil, ErrUnsupported
	}
}

// regex is the builtin validation function that checks
// whether the string variable matches a regular expression
func regex(t reflect.Type, param string) (validationFunc, error) {
	re, err := regexp.Compile(param)
	if err != nil {
		return nil, ErrBadParameter
	}
	if t != reflect.TypeOf("") {
		return nil, ErrUnsupported
	}
	return func(v reflect.Value, state *validateState) {
		if !re.MatchString(v.String()) {
			state.error(ErrRegexp)
		}
	}, nil
}
