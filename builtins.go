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
	"reflect"
	"regexp"
	"strconv"
)

const (
	errZeroValue = "zero value"
	errMin       = "less than min"
	errMax       = "greater than max"
	errLen       = "invalid length"
	errRegexp    = "regular expression mismatch"
	errRequired  = "required value"
)

var errUnsupportedType = errors.New("unsupported type")
var errBadParameter = errors.New("bad parameter")

// nonzero tests whether a variable value non-zero
// as defined by the golang spec.
func nonzero(t reflect.Type, param string) (ValidationFunc, error) {
	check := func(ok bool, r *ErrorReporter) {
		if !ok {
			r.Errorf(errZeroValue)
		}
	}
	switch t.Kind() {
	case reflect.String:
		return func(v reflect.Value, r *ErrorReporter) {
			check(len(v.String()) != 0, r)
		}, nil
	case reflect.Ptr, reflect.Interface:
		return func(v reflect.Value, r *ErrorReporter) {
			check(!v.IsNil(), r)
		}, nil
	case reflect.Slice, reflect.Map, reflect.Array:
		return func(v reflect.Value, r *ErrorReporter) {
			check(v.Len() != 0, r)
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(v reflect.Value, r *ErrorReporter) {
			check(v.Int() != 0, r)
		}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return func(v reflect.Value, r *ErrorReporter) {
			check(v.Uint() != 0, r)
		}, nil
	case reflect.Float32, reflect.Float64:
		return func(v reflect.Value, r *ErrorReporter) {
			// TODO this preserves old behavior but is arguably
			// incorrect - it should probably error if the value is NaN.
			check(v.Float() != 0, r)
		}, nil
	case reflect.Bool:
		return func(v reflect.Value, r *ErrorReporter) {
			check(v.Bool(), r)
		}, nil
	case reflect.Struct:
		return okValidation, nil
	}
	return nil, errUnsupportedType
}

// length tests whether a variable's length is equal to a given
// value. For strings it tests the number of characters whereas
// for maps and slices it tests the number of items.
func length(t reflect.Type, param string) (ValidationFunc, error) {
	switch t.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if int64(v.Len()) != p {
				r.Errorf(errLen)
			}
		}, nil
	}
	return nil, errUnsupportedType
}

// min tests whether a variable value is larger or equal to a given
// number. For number types, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func min(t reflect.Type, param string) (ValidationFunc, error) {
	switch t.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if int64(v.Len()) < p {
				r.Errorf(errMin)
			}
		}, nil
	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if v.Float() < p {
				r.Errorf(errMin)
			}
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if v.Int() < p {
				r.Errorf(errMin)
			}
		}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := strconv.ParseUint(param, 0, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if v.Uint() < p {
				r.Errorf(errMin)
			}
		}, nil
	default:
		return nil, errUnsupportedType
	}
}

// max tests whether a variable value is lesser than a given
// value. For numbers, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func max(t reflect.Type, param string) (ValidationFunc, error) {
	switch t.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if int64(v.Len()) > p {
				r.Errorf(errMax)
			}
		}, nil
	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if v.Float() > p {
				r.Errorf(errMax)
			}
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if v.Int() > p {
				r.Errorf(errMax)
			}
		}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := strconv.ParseUint(param, 0, 64)
		if err != nil {
			return nil, errBadParameter
		}
		return func(v reflect.Value, r *ErrorReporter) {
			if v.Uint() > p {
				r.Errorf(errMax)
			}
		}, nil
	default:
		return nil, errUnsupportedType
	}
}

// regex is the builtin validation function that checks
// whether the string variable matches a regular expression
func regex(t reflect.Type, param string) (ValidationFunc, error) {
	re, err := regexp.Compile(param)
	if err != nil {
		return nil, errBadParameter
	}
	if t != reflect.TypeOf("") {
		return nil, errUnsupportedType
	}
	return func(v reflect.Value, r *ErrorReporter) {
		if !re.MatchString(v.String()) {
			r.Errorf(errRegexp)
		}
	}, nil
}

// required validates the value is not nil for a field, that is, a
// pointer or an interface, any other case is a valid one as zero
// value from Go spec.
func required(t reflect.Type, param string) (ValidationFunc, error) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice:
		return func(v reflect.Value, r *ErrorReporter) {
			if v.IsNil() {
				r.Errorf(errRequired)
			}
		}, nil
	case reflect.String, reflect.Array, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool, reflect.Struct:
		return okValidation, nil
	default:
		return nil, errUnsupportedType
	}
}

var uuidRE = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// uuid validates if a string represents a valid UUID (RFC 4122)
func uuid(t reflect.Type, param string) (ValidationFunc, error) {
	if t.Kind() != reflect.String {
		return nil, errUnsupportedType
	}
	return func(v reflect.Value, r *ErrorReporter) {
		if !uuidRE.MatchString(v.String()) {
			r.Errorf(errRegexp)
		}
	}, nil
}

// latitude validates that a field is a latitude
func latitude(t reflect.Type, param string) (ValidationFunc, error) {
	validateLatitude := func(f float64, r *ErrorReporter) {
		if f < -90 || f > 90 {
			r.Errorf("%g is not a valid latitude", f)
		}
	}

	switch t.Kind() {
	case reflect.Float64:
		return func(v reflect.Value, r *ErrorReporter) {
			validateLatitude(v.Float(), r)
		}, nil
	case reflect.String:
		return func(v reflect.Value, r *ErrorReporter) {
			s := v.String()
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				r.Errorf("%q is not a valid latitude", s)
			}
			validateLatitude(f, r)
		}, nil
	default:
		return nil, errUnsupportedType
	}
}

// longitude validates that a field is a longitude
func longitude(t reflect.Type, param string) (ValidationFunc, error) {
	validateLongitude := func(f float64, r *ErrorReporter) {
		if f < -180 || f > 180 {
			r.Errorf("%g is not a valid longitude", f)
		}
	}
	switch t.Kind() {
	case reflect.Float64:
		return func(v reflect.Value, r *ErrorReporter) {
			validateLongitude(v.Float(), r)
		}, nil
	case reflect.String:
		return func(v reflect.Value, r *ErrorReporter) {
			s := v.String()
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				r.Errorf("%q is not a valid longitude", s)
			}
			validateLongitude(f, r)
		}, nil
	default:
		return nil, errUnsupportedType
	}
}
