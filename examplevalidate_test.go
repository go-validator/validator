// Package walidator implements value validations
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

package walidator_test

import (
	"encoding/json"
	"fmt"

	"github.com/heetch/walidator/v4"
)

// This example demonstrates a custom function to process template text.
// It installs the strings.Title function and uses it to
// Make Title Text Look Good In Our Template's Output.
func ExampleValidate() {
	// First create a struct to be validated
	// according to the validate tags.
	type ValidateExample struct {
		Name        string `validate:"nonzero"`
		Description string
		Age         int    `validate:"min=18"`
		Email       string `validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
		Address     struct {
			Street string `validate:"nonzero"`
			City   string `validate:"nonzero"`
		}
	}

	// Fill in some values
	ve := ValidateExample{
		Name:        "Joe Doe", // valid as it's nonzero
		Description: "",        // valid no validation tag exists
		Age:         17,        // invalid as age is less than required 18
	}
	// invalid as Email won't match the regular expression
	ve.Email = "@not.a.valid.email"
	ve.Address.City = "Some City" // valid
	ve.Address.Street = ""        // invalid

	err := walidator.Validate(ve)
	if err == nil {
		fmt.Println("Values are valid.")
		return
	}

	fmt.Println("Invalid with errors:")
	data, _ := json.MarshalIndent(err, "\t", "\t")
	fmt.Println("\t" + string(data))
	// Output:
	//Invalid with errors:
	//	{
	//		"Address.Street": [
	//			{
	//				"kind": "nonzero",
	//				"msg": "zero value"
	//			}
	//		],
	//		"Age": [
	//			{
	//				"kind": "min",
	//				"msg": "less than min"
	//			}
	//		],
	//		"Email": [
	//			{
	//				"kind": "regexp",
	//				"msg": "regular expression mismatch"
	//			}
	//		]
	//	}
}

// This example shows you how to change the tag name
func ExampleWithTag() {
	type T struct {
		A int `foo:"nonzero" bar:"min=10"`
	}
	t := T{5}
	v := walidator.New().WithTag("foo")
	err := v.Validate(t)
	fmt.Printf("foo --> valid: %v, errs: %v\n", err == nil, err)

	v = walidator.New().WithTag("bar")
	err = v.Validate(t)
	fmt.Printf("bar --> valid: %v, errs: %v\n", err == nil, err)

	// Output:
	// foo --> valid: true, errs: <nil>
	// bar --> valid: false, errs: validation error: A: less than min
}
