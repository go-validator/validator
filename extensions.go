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

package validator

import "reflect"

// required validates the value is not nil for a field, that is, a
// pointer or an interface, any other case is a valid one as zero
// value from Go spec
func required(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	var valid bool
	switch st.Kind() {
	case reflect.Ptr, reflect.Interface:
		valid = !st.IsNil()
	case reflect.Invalid:
		valid = false
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool, reflect.Struct:
		valid = true
	default:
		return ErrUnsupported
	}
	if valid {
		return nil
	}
	return ErrRequired
}

// uuid validates if a string represents a valid UUID (RFC 4122)
func uuid(v interface{}, param string) error {
	uuidRE := "(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	return regex(v, uuidRE)
}
