// Package validator_test test value validations
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

package validator_test

import (
	. "gopkg.in/check.v1"

	validator "github.com/heetch/walidator"
)

type ExtensionSuite struct{}

var _ = Suite(&ExtensionSuite{})

func (es *ExtensionSuite) TestUUIDOK(c *C) {
	cases := []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"0FCE98AC-1326-4C79-8EBC-94908DA8B034",
	}
	for _, s := range cases {
		err := validator.Valid(s, "uuid")
		c.Assert(err, IsNil)
	}
}

func (es *ExtensionSuite) TestUUIDNOK(c *C) {
	cases := []string{
		"1234",
		"0VCE98AC-1326-4C79-8EBC-94908DA8B034",
	}
	for _, s := range cases {
		err := validator.Valid(s, "uuid")
		c.Assert(err, NotNil)
		errs, ok := err.(validator.ErrorArray)
		c.Assert(ok, Equals, true)
		c.Assert(errs, HasLen, 1)
		c.Assert(errs, HasError, validator.ErrRegexp)
	}
}

type Mer interface {
	M()
}

type T1 struct{}

func (t *T1) M() {}

type T2 struct {
	Mer Mer
}

func (es *ExtensionSuite) TestRequiredOK(c *C) {
	a := []int{1, 2, 3}
	cases := []interface{}{
		"string",
		a[1:],
		a,
		map[string]int{"a": 1, "b": 2},
		12,
		2.1,
		T1{},
		struct{ Foo int }{23},
	}
	for _, s := range cases {
		err := validator.Valid(s, "required")
		c.Assert(err, IsNil)
	}
}

func (es *ExtensionSuite) TestRequiredNOK(c *C) {
	var ptr *uint
	var t1 *T1
	t2 := T2{}
	cases := []interface{}{
		ptr,
		t1,
		nil,
		t2.Mer,
	}
	for _, s := range cases {
		err := validator.Valid(s, "required")
		c.Assert(err, NotNil)
		errs, ok := err.(validator.ErrorArray)
		c.Assert(ok, Equals, true)
		c.Assert(errs, HasLen, 1)
		c.Assert(errs, HasError, validator.ErrRequired)
	}
}
