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

package walidator_test

import (
	"fmt"

	"github.com/heetch/walidator"
	. "gopkg.in/check.v1"
)

type ExtensionSuite struct{}

var _ = Suite(&ExtensionSuite{})

func (es *ExtensionSuite) TestUUIDOK(c *C) {
	cases := []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"0FCE98AC-1326-4C79-8EBC-94908DA8B034",
	}
	for _, s := range cases {
		err := walidator.Valid(s, "uuid")
		c.Assert(err, IsNil)
	}
}

func (es *ExtensionSuite) TestUUIDNOK(c *C) {
	cases := []string{
		"1234",
		"0VCE98AC-1326-4C79-8EBC-94908DA8B034",
	}
	for _, s := range cases {
		err := walidator.Valid(s, "uuid")
		c.Assert(err, NotNil)
		errs, ok := err.(walidator.ErrorArray)
		c.Assert(ok, Equals, true)
		c.Assert(errs, HasLen, 1)
		c.Assert(errs, HasError, walidator.ErrRegexp)
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
		err := walidator.Valid(s, "required")
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
		err := walidator.Valid(s, "required")
		c.Assert(err, NotNil)
		errs, ok := err.(walidator.ErrorArray)
		c.Assert(ok, Equals, true)
		c.Assert(errs, HasLen, 1)
		c.Assert(errs, HasError, walidator.ErrRequired)
	}
}

func (es *ExtensionSuite) TestLatitudeOK(c *C) {
	v := 1.1
	s := "-12.1"
	cases := []interface{}{
		21.0,
		0.0,
		-10.212,
		"12.0",
		"0.0",
		"-1.22",
		&v,
		&s,
	}
	for _, l := range cases {
		err := walidator.Valid(l, "latitude")
		c.Assert(err, IsNil)
	}
}

func (es *ExtensionSuite) TestLatitudeNOK(c *C) {
	v := 220.21
	s := "-2121.1"
	cases := []interface{}{
		210.0,
		-1000.212,
		"1220",
		"-190.2",
		&v,
		&s,
	}
	for _, l := range cases {
		err := walidator.Valid(l, "latitude")
		c.Assert(err, NotNil)
		c.Assert(err, FitsTypeOf, walidator.ErrorArray{})
		errs, ok := err.(walidator.ErrorArray)
		c.Assert(ok, Equals, true)
		c.Assert(errs, HasLen, 1)
		switch loc := l.(type) {
		case *float64:
			c.Assert(errs[0].Error(), Equals, fmt.Sprintf("%v is not a valid latitude", *loc))
		case *string:
			c.Assert(errs[0].Error(), Equals, fmt.Sprintf("%v is not a valid latitude", *loc))
		default:
			c.Assert(errs[0].Error(), Equals, fmt.Sprintf("%v is not a valid latitude", loc))
		}
	}
}

func (es *ExtensionSuite) TestLongitudeOK(c *C) {
	v := 1.1
	s := "-12.1"
	cases := []interface{}{
		21.0,
		0.0,
		-10.212,
		"12.0",
		"0.0",
		"-1.22",
		&v,
		&s,
	}
	for _, l := range cases {
		err := walidator.Valid(l, "latitude")
		c.Assert(err, IsNil)
	}
}

func (es *ExtensionSuite) TestLongitudeNOK(c *C) {
	v := 220.21
	s := "-2121.1"
	cases := []interface{}{
		210.0,
		-1000.212,
		"1220",
		"-190.2",
		&v,
		&s,
	}
	for _, l := range cases {
		err := walidator.Valid(l, "longitude")
		c.Assert(err, NotNil)
		errs, ok := err.(walidator.ErrorArray)
		c.Assert(ok, Equals, true)
		c.Assert(errs, HasLen, 1)
		switch loc := l.(type) {
		case *float64:
			c.Assert(errs[0].Error(), Equals, fmt.Sprintf("%v is not a valid longitude", *loc))
		case *string:
			c.Assert(errs[0].Error(), Equals, fmt.Sprintf("%v is not a valid longitude", *loc))
		default:
			c.Assert(errs[0].Error(), Equals, fmt.Sprintf("%v is not a valid longitude", loc))
		}
	}
}
