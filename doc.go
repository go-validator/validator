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

/*
Package validator implements value validations based on struct tags.

In code it is often necessary to validate that a given value is valid before
using it for something. A typical example might be something like this.

	if age < 18 {
		return error.New("age cannot be under 18")
	}

This is a simple enough example, but it can get significantly more complex,
especially when dealing with structs.

	l := len(strings.Trim(s.Username))
	if l < 3 || l > 40  || !regexp.MatchString("^[a-zA-Z]$", s.Username) ||	s.Age < 18 || s.Password {
		return errors.New("Invalid request")
	}

You get the idea. Package validator allows one to define valid values as
struct tags when defining a new struct type.

	type NewUserRequest struct {
		Username string `validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
		Name string     `validate:"nonzero"`
		Age int         `validate:"min=18"`
		Password string `validate:"min=8"`
	}

Then validating a variable of type NewUserRequest becomes trivial.

	r := NewUserRequest{Username: "something", ...}
	if errs := walidator.Validate(r); errs != nil {
		// do something
	}

Builtin validator functions

Here is the list of validator functions builtin in the package.

	len
		For numeric numbers, len will simply make sure that the value is
		equal to the parameter given. For strings, it checks that
		the string length is exactly that number of characters. For slices,
		arrays, and maps, validates the number of items. (Usage: len=10)

	max
		For numeric numbers, max will simply make sure that the value is
		lesser or equal to the parameter given. For strings, it checks that
		the string length is at most that number of characters. For slices,
		arrays, and maps, validates the number of items. (Usage: max=10)

	min
		For numeric numbers, min will simply make sure that the value is
		greater or equal to the parameter given. For strings, it checks that
		the string length is at least that number of characters. For slices,
		arrays, and maps, validates the number of items. (Usage: min=10)

	nonzero
		This validates that the value is not zero. The appropriate zero value
		is given by the Go spec (e.g. for int it's 0, for string it's "", for
		pointers is nil, etc.) Usage: nonzero

	required
		For pointer and interface types, this validates that the value
		is non-nil.

	regexp
		Only valid for string types, it will validate that the value matches
		the regular expression provided as parameter. (Usage: regexp=^a.*b$)

	uuid
		This validates that the string value is a valid RFC 4122 UUID.

	latitude
		This validates that a float64 or string value contains a valid
		geographical latitude value (between -90 and 90 degrees)

	longitude
		This validates that a float64 or string value contains a valid
		geographical longitude value (between -180 and 180 degrees).


Note that there are no tests to prevent conflicting validator parameters. For
instance, these fields will never be valid.

	...
	A int     `validate:"max=0,min=1"`
	B string  `validate:"len=10,regexp=^$"
	...

Multiple validators

You may often need to have a different set of validation
rules for different situations. The global default validator
cannot be changed (to avoid global name clashes) but
you can create custom Validator instances by calling
walidator.New.

You can add your own custom validation functions to
a Validator instance and you can also change the struct
field tag that's recognized from the default "validate" tag.

Custom validation functions

It is possible to add your own custom validation functions with AddValidation.
Note that the validation function is passed the actual type that will be used
when validating. This enables it to create a specialized function
for the type and parameter, and to return an early error when the type is known
to be incorrect.

	// notZZ defines a validation function that checks that
	// a string value is not "ZZ"
	func notZZ(t reflect.Type, param string) (walidator.ValidationFunc, error) {
		if t.Kind() != reflect.String {
			return nil, fmt.Errorf("unsupported type")
		}
		return func(v reflect.Value, r *walidator.ErrorReporter) {
			if v.String() == "ZZ" {
				r.Errorf("value cannot be ZZ")
			}
		}, nil
	}

Then one needs to add it to the list of validation funcs and give it a "tag" name.

	v := walidator.New()
	v.AddValidation("notzz", notZZ)

Then it is possible to use the notzz validation tag. This will print
"validation error: A: value cannot be ZZ"

	type T struct {
		A string  `validate:"nonzero,notzz"`
	}
	t := T{"ZZ"}
	if err := v.Validate(t); err != nil {
		fmt.Printf("validation error: %v\n", err)
	}

If you wish to parameterize the tag, you can use the parameter
argument to the validator, which holds the string after the "="
in the struct tag.

	// notSomething defines a validation function that checks that
	// a string value is not equal to its parameter.
	func notSomething(t reflect.Type, param string) (walidator.ValidationFunc, error) {
		if t.Kind() != reflect.String {
			return nil, fmt.Errorf("unsupported type")
		}
		return func(v reflect.Value, r *walidator.ErrorReporter) {
			if v.String() == param {
				r.Errorf("value cannot be %q", param)
			}
		}, nil
	}

And then the code below should print "validation error: A: value cannot be ABC".

	v.AddValidation("notsomething", notSomething)
	type T struct {
		A string  `validate:"notsomething=ABC"`
	}
	t := T{"ABC"}
	if err := v.Validate(t); err != nil {
		fmt.Printf("validation error: %v\n", err)
	}

Using an unknown validation func in a field tag will cause Validate to return
a validation error.

Custom tag name

There might be a reason not to use the tag 'validate' (maybe due to
a conflict with a different package). In this case, you can
create a new Validator instance that will use a different tag:

	v := walidator.New().WithTag("valid")

Then.

	Type T struct {
		A int    `valid:"min=8, max=10"`
		B string `valid:"nonzero"`
	}
*/
package walidator
