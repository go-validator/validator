// Package validator implements value validations
//
// Copyright 2018 Heetch
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
)

// required validates the value is not nil for a field, that is, a
// pointer or an interface, any other case is a valid one as zero
// value from Go spec
func required(t reflect.Type, param string) (validationFunc, error) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface:
		return func(v reflect.Value, state *validateState) {
			if v.IsNil() {
				state.error(ErrRequired)
			}
		}, nil
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool, reflect.Struct:
		return okValidation, nil
	default:
		return nil, ErrUnsupported
	}
}

var uuidRE = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// uuid validates if a string represents a valid UUID (RFC 4122)
func uuid(t reflect.Type, param string) (validationFunc, error) {
	if t != reflect.TypeOf("") {
		return nil, ErrUnsupported
	}
	return func(v reflect.Value, state *validateState) {
		if !uuidRE.MatchString(v.Interface().(string)) {
			state.error(ErrRegexp)
		}
	}, nil
}

// latitude validates that a field is a latitude
func latitude(t reflect.Type, param string) (validationFunc, error) {
	validateLatitude := func(f float64, state *validateState) {
		if f < -90 || f > 90 {
			state.error(TextErr{Err: fmt.Errorf("%g is not a valid latitude", f)})
		}
	}

	switch t.Kind() {
	case reflect.Float64:
		return func(v reflect.Value, state *validateState) {
			validateLatitude(v.Float(), state)
		}, nil
	case reflect.String:
		return func(v reflect.Value, state *validateState) {
			s := v.String()
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				state.error(TextErr{Err: fmt.Errorf("%g is not a valid latitude", f)})
			}
			validateLatitude(f, state)
		}, nil
	default:
		return nil, ErrUnsupported
	}
}

// longitude validates that a field is a longitude
func longitude(t reflect.Type, param string) (validationFunc, error) {
	validateLongitude := func(f float64, state *validateState) {
		if f < -180 || f > 180 {
			state.error(TextErr{Err: fmt.Errorf("%g is not a valid longitude", f)})
		}
	}
	switch t.Kind() {
	case reflect.Float64:
		return func(v reflect.Value, state *validateState) {
			validateLongitude(v.Float(), state)
		}, nil
	case reflect.String:
		return func(v reflect.Value, state *validateState) {
			s := v.String()
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				state.error(TextErr{Err: fmt.Errorf("%g is not a valid latitude", f)})
			}
			validateLongitude(f, state)
		}, nil
	default:
		return nil, ErrUnsupported
	}
}
