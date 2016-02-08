// Package validator implements value validations
//
// Copyright (C) 2014-2016 Roberto Teixeira <robteix@robteix.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package validator_test

import (
	"fmt"

	"gopkg.in/validator.v2"
)

// This example demonstrates a custom function to process template text.
// It installs the strings.Title function and uses it to
// Make Title Text Look Good In Our Template's Output.
func ExampleValidate() {
	// First create a struct to be validated
	// according to the validator tags.
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

	valid, err := validator.Validate(ve)
	if valid {
		fmt.Println("Values are valid.")
	} else {
		// See if Address was empty
		if validator.IsZeroValue(err, "Address.Street") {
			fmt.Println("Street cannot be empty.")
		}

		fmt.Println("Errors found:")
		for _, f := range validator.ErrorFields(err) {
			fmt.Printf("\t - %s (%v)\n", f, validator.Errors(err, f))
		}
	}

	// Output:
	// Street cannot be empty.
	// Errors found:
	//	 - Age ([less than min])
	//	 - Email ([regular expression mismatch])
	//	 - Address.Street ([zero value])

}

// This example shows how to use the Valid helper
// function to validator any number of values
func ExampleValid() {
	valid, errs := validator.Valid(42, "min=10,max=100,nonzero")
	fmt.Printf("42: valid=%v, errs=%v\n", valid, errs)

	var ptr *int
	if valid, _ := validator.Valid(ptr, "nonzero"); !valid {
		fmt.Println("ptr: Invalid nil pointer.")
	}

	valid, _ = validator.Valid("ABBA", "regexp=[ABC]*")
	fmt.Printf("ABBA: valid=%v\n", valid)

	// Output:
	// 42: valid=true, errs=[]
	// ptr: Invalid nil pointer.
	// ABBA: valid=true
}

// This example shows you how to change the tag name
func ExampleSetTag() {
	type T struct {
		A int `foo:"nonzero" bar:"min=10"`
	}
	t := T{5}
	v := validator.NewValidator()
	v.SetTag("foo")
	valid, err := v.Validate(t)
	fmt.Printf("foo --> valid: %v, err: %v\n", valid, err)
	v.SetTag("bar")
	valid, err = v.Validate(t)
	fmt.Printf("bar --> valid: %v, err: %v\n", valid, err)

	// Output:
	// foo --> valid: true, err: <nil>
	// bar --> valid: false, err: A has error less than min
}

// This example shows you how to change the tag name
func ExampleWithTag() {
	type T struct {
		A int `foo:"nonzero" bar:"min=10"`
	}
	t := T{5}
	valid, errs := validator.WithTag("foo").Validate(t)
	fmt.Printf("foo --> valid: %v, errs: %v\n", valid, errs)
	valid, errs = validator.WithTag("bar").Validate(t)
	fmt.Printf("bar --> valid: %v, errs: %v\n", valid, errs)

	// Output:
	// foo --> valid: true, errs: <nil>
	// bar --> valid: false, errs: A has error less than min
}
